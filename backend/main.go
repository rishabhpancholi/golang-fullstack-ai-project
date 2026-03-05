package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"project/config"
	db "project/database"
	"project/routes"
	utilities "project/utils"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v", err.Error())
	}

	pool, err := db.Connect(cfg.DataBaseURL)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err.Error())
	}

	cld, err := utilities.CloudinarySetup(cfg.CloudinaryURL)
	if err != nil {
		log.Fatalf("Error setting up cloudinary: %v", err.Error())
	}

	gin.SetMode(gin.ReleaseMode)

	app := gin.Default()

	app.MaxMultipartMemory = 10 << 20

	app.SetTrustedProxies(nil)

	routes.RegisterAuthRoutes(app, pool, cld, &cfg.JWTSecret)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: app,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting the server: %v", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("Closing database connection...")

	pool.Close()

	log.Println("Shutting down server...")

	if err := srv.Shutdown(ctx); err != nil {
		log.Println("Server shutdown: ", err)
	}

	log.Println("Server exited")
}
