package models

import (
	"mime/multipart"
	"time"
)

// Register user request model
type RegisterUserInput struct {
	Name           string  `json:"name" form:"name" binding:"required"`
	Email          string  `json:"email" form:"email" binding:"required,email"`
	Password       string  `json:"password" form:"password" binding:"required"`
	PhoneNumber    string  `json:"phone_number" form:"phone_number" binding:"required"`
	Bio            *string `json:"bio" form:"bio"`
	ResumeFile     *multipart.FileHeader `json:"resume_file" form:"resume_file"`
	Resume         *string `json:"resume" form:"resume"`
	ResumePublicId *string `json:"resume_public_id" form:"resume_public_id"`
}

type RegisteredUser struct {
	UserId      int       `json:"user_id" db:"user_id"`
	Name        string    `json:"name" db:"name"`
	Email       string    `json:"email" db:"email"`
	PhoneNumber string    `json:"phone_number" db:"phone_number"`
	Role        string    `json:"role" db:"role"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}
