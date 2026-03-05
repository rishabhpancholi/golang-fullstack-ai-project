package handlers

import (
	"net/http"
	"path/filepath"
	"project/models"
	"project/repositories"
	utilities "project/utils"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Handler for register jobseeker requests
func RegisterJobseekerHandler(pool *pgxpool.Pool, cld *cloudinary.Cloudinary, jwtSecret *string) gin.HandlerFunc {
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

		uploadResultURL, uploadResultPublicId, err := utilities.CloudinaryUpload(resume, cld)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error uploading resume, please try again later",
			})
			return
		}

		user_input.Resume = &uploadResultURL
		user_input.ResumePublicId = &uploadResultPublicId

		existing_id, err := repositories.GetUserWithEmail(pool, user_input.Email)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error, please try again later",
			})
			return
		}

		if existing_id != 0 {
			ctx.JSON(http.StatusConflict, gin.H{
				"error": "User with this email already exists",
			})
			return
		}

		var registered_user *models.RegisteredUser

		registered_user, err = repositories.RegisterJobseeker(pool, user_input)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusCreated, registered_user)
	}
}

// Handler for register jobseeker requests
func RegisterRecruiterHandler(pool *pgxpool.Pool, cld *cloudinary.Cloudinary, jwtSecret *string) gin.HandlerFunc {
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

		existing_id, err := repositories.GetUserWithEmail(pool, user_input.Email)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error, please try again later",
			})
			return
		}

		if existing_id != 0 {
			ctx.JSON(http.StatusConflict, gin.H{
				"error": "User with this email already exists",
			})
			return
		}

		var registered_user *models.RegisteredUser

		registered_user, err = repositories.RegisterRecruiter(pool, user_input)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error, please try again later",
			})
			return
		}

		ctx.JSON(http.StatusCreated, registered_user)
	}
}
