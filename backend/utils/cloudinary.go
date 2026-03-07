package utilities

import (
	"context"
	"mime/multipart"
	"project/config"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

// Sets up cloudinary
func CloudinarySetup(cloudinaryURL string) (*cloudinary.Cloudinary, error) {
	cld, err := cloudinary.NewFromURL(cloudinaryURL)
	if err != nil {
		return nil, err
	}

	return cld, nil
}

func CloudinaryUpload(file multipart.File, app *config.App) (string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	uploadResult, err := app.Cld.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder: "resumes",
	})
	if err != nil {
		return "", "", err
	}

	return uploadResult.SecureURL, uploadResult.PublicID, nil
}
