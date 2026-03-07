package handlers

import (
	"net/http"
	"path/filepath"
	"project/config"
	"project/models"
	"project/repositories"
	utilities "project/utils"

	"github.com/gin-gonic/gin"
)

// Handler for register jobseeker requests
func RegisterJobseekerHandler(app *config.App) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var user_input *models.RegisterUserInput

		if err := ctx.ShouldBind(&user_input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Please fill the details correctly",
			})
			return
		}

		if len(user_input.PhoneNumber) != 10 {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid phone number",
			})
			return
		}

		existing_id, _, err := repositories.GetUserWithEmail(app, user_input.Email)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error, please try again later",
			})
			return
		}

		if existing_id != "" {
			ctx.JSON(http.StatusConflict, gin.H{
				"error": "User with this email already exists",
			})
			return
		}

		resumeFileExt := filepath.Ext(user_input.ResumeFile.Filename)
		if resumeFileExt != ".pdf" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Please upload pdf files only",
			})
			return
		}

		resume, err := user_input.ResumeFile.Open()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Could not open resume file, please try again later",
			})
			return
		}

		uploadResultURL, uploadResultPublicId, err := utilities.CloudinaryUpload(resume, app)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error uploading resume, please try again later",
			})
			return
		}

		user_input.Resume = &uploadResultURL
		user_input.ResumePublicId = &uploadResultPublicId

		var registered_user *models.RegisteredUser

		registered_user, token, err := repositories.RegisterJobseeker(app, user_input)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"user":  registered_user,
			"token": token,
		})
	}
}

// Handler for register jobseeker requests
func RegisterRecruiterHandler(app *config.App) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var user_input *models.RegisterUserInput

		if err := ctx.ShouldBind(&user_input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Please fill the details correctly",
			})
			return
		}

		if len(user_input.PhoneNumber) != 10 {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid phone number",
			})
			return
		}

		existing_id, _, err := repositories.GetUserWithEmail(app, user_input.Email)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		if existing_id != "" {
			ctx.JSON(http.StatusConflict, gin.H{
				"error": "User with this email already exists",
			})
			return
		}

		var registered_user *models.RegisteredUser

		registered_user, token, err := repositories.RegisterRecruiter(app, user_input)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"user":  registered_user,
			"token": token,
		})
	}
}

func UserLoginHandler(app *config.App) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input_credentials *models.UserLoginInput

		if err := ctx.ShouldBind(&input_credentials); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Please give proper login credentials",
			})
			return
		}

		existing_id, existingPassword, err := repositories.GetUserWithEmail(app, input_credentials.Email)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error, please try again later",
			})
			return
		}

		if existing_id == "" {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "User with this email does not exist",
			})
			return
		}

		if err := utilities.VerifyPassword(existingPassword, input_credentials.Password); err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid credentials",
			})
			return
		}

		token, err := repositories.LoginUser(app, input_credentials)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error generating token, please try again later",
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"token": token,
		})
	}
}
