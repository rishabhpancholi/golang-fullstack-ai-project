package routes

import (
	"project/config"
	"project/handlers"
)

// Registers the auth routes
func RegisterAuthRoutes(app *config.App) {
	auth := app.Handler.Group("/auth")

	auth.POST("/register-jobseeker", handlers.RegisterJobseekerHandler(app))
	auth.POST("/register-recruiter", handlers.RegisterRecruiterHandler(app))
	auth.POST("/login-user", handlers.UserLoginHandler(app))
}
