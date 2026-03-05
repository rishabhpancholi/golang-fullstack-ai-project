package routes

import (
	"project/handlers"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Registers the auth routes
func RegisterAuthRoutes(app *gin.Engine, pool *pgxpool.Pool, cld *cloudinary.Cloudinary, jwtSecret *string) {
	auth := app.Group("/auth")

	auth.POST("/register-jobseeker", handlers.RegisterJobseekerHandler(pool, cld, jwtSecret))
	auth.POST("/register-recruiter", handlers.RegisterRecruiterHandler(pool, cld, jwtSecret))
}
