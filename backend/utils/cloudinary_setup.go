package utilities

import "github.com/cloudinary/cloudinary-go/v2"

// Sets up cloudinary
func CloudinarySetup(cloudinaryURL string) (*cloudinary.Cloudinary, error) {
	cld, err := cloudinary.NewFromURL(cloudinaryURL)
	if err != nil {
		return nil, err
	}

	return cld, nil
}
