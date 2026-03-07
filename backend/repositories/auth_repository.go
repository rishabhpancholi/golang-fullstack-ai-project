package repositories

import (
	"context"
	"errors"
	"project/config"
	"project/models"
	utilities "project/utils"
	"time"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

// Gets user from db with the provided email
func GetUserWithEmail(app *config.App, email string) (string, string, error) {
	var userId string
	var userPassword string

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sqlStatement := `
	SELECT user_id, password from users
	WHERE email = $1;
	`
	if err := app.Pool.QueryRow(ctx, sqlStatement, email).Scan(
		&userId,
		&userPassword,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", "", nil
		}

		return "", "", err
	}

	return userId, userPassword, nil
}

// Registers a new jobseeker in db
func RegisterJobseeker(app *config.App, user_input *models.RegisterUserInput) (*models.RegisteredUser, string, error) {
	var registered_user models.RegisteredUser

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sqlStatement := `
	INSERT INTO users (name, email, password, phone_number, role, bio, resume, resume_public_id)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING user_id, name, email, phone_number, role, created_at
	`
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user_input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	if err := app.Pool.QueryRow(
		ctx,
		sqlStatement,
		user_input.Name,
		user_input.Email,
		string(hashedPassword),
		user_input.PhoneNumber,
		"jobseeker",
		user_input.Bio,
		user_input.Resume,
		user_input.ResumePublicId,
	).Scan(
		&registered_user.UserId,
		&registered_user.Name,
		&registered_user.Email,
		&registered_user.PhoneNumber,
		&registered_user.Role,
		&registered_user.CreatedAt,
	); err != nil {
		return nil, "", err
	}

	token, err := utilities.GenerateJWTToken(app, &registered_user)
	if err != nil {
		return nil, "", err
	}

	return &registered_user, token, nil
}

// Registers a new recruiter in db
func RegisterRecruiter(app *config.App, user_input *models.RegisterUserInput) (*models.RegisteredUser, string, error) {
	var registered_user models.RegisteredUser

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sqlStatement := `
	INSERT INTO users (name, email, password, phone_number, role, bio)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING user_id, name, email, phone_number, role, created_at
	`
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user_input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	if err := app.Pool.QueryRow(
		ctx,
		sqlStatement,
		user_input.Name,
		user_input.Email,
		string(hashedPassword),
		user_input.PhoneNumber,
		"recruiter",
		user_input.Bio,
	).Scan(
		&registered_user.UserId,
		&registered_user.Name,
		&registered_user.Email,
		&registered_user.PhoneNumber,
		&registered_user.Role,
		&registered_user.CreatedAt,
	); err != nil {
		return nil, "", err
	}

	token, err := utilities.GenerateJWTToken(app, &registered_user)
	if err != nil {
		return nil, "", err
	}

	return &registered_user, token, nil
}

func LoginUser(app *config.App, input_credentials *models.UserLoginInput) (string, error) {
	var registered_user models.RegisteredUser

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sqlStatement := `
	SELECT user_id, email FROM users 
	WHERE email = $1;
	`

	if err := app.Pool.QueryRow(
		ctx,
		sqlStatement,
		input_credentials.Email,
	).Scan(
		&registered_user.UserId,
		&registered_user.Email,
	); err != nil {
		return "", err
	}

	token, err := utilities.GenerateJWTToken(app, &registered_user)
	if err != nil {
		return "", err
	}

	return token, nil
}
