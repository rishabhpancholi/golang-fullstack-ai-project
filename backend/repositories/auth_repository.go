package repositories

import (
	"context"
	"errors"
	"project/models"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

// Gets user from db with the provided email
func GetUserWithEmail(pool *pgxpool.Pool, email string) (int, error) {
	var userId int

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sqlStatement := `
	SELECT user_id from users
	WHERE email = $1;
	`
	if err := pool.QueryRow(ctx, sqlStatement, email).Scan(
		&userId,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, nil
		}

		return 0, err
	}

	return userId, nil
}

// Registers a new jobseeker in db
func RegisterJobseeker(pool *pgxpool.Pool, user_input *models.RegisterUserInput) (*models.RegisteredUser, error) {
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
		return nil, err
	}

	if err := pool.QueryRow(
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
		return nil, err
	}

	return &registered_user, nil
}

// Registers a new recruiter in db
func RegisterRecruiter(pool *pgxpool.Pool, user_input *models.RegisterUserInput) (*models.RegisteredUser, error) {
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
		return nil, err
	}

	if err := pool.QueryRow(
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
		return nil, err
	}

	return &registered_user, nil
}
