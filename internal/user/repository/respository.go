package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/pi-prakhar/go-gcp-pi-app/internal/auth/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user models.GoogleUser) error
	GetUserByEmail(ctx context.Context, email string) (*models.GoogleUser, error)
}

type GCPPostgresqlRepository struct {
	DB *sql.DB
}

// Create a new GoogleUser
func (r *GCPPostgresqlRepository) CreateUser(ctx context.Context, user models.GoogleUser) error {
	query := `
        INSERT INTO users (id, email, verified_email, name, given_name, family_name, picture, locale)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.DB.ExecContext(ctx, query, user.ID, user.Email, user.VerifiedMail, user.Name, user.GivenName, user.FamilyName, user.Picture, user.Locale)
	if err != nil {
		return err
	}
	return nil
}

// Get a GoogleUser by email
func (r *GCPPostgresqlRepository) GetUserByEmail(ctx context.Context, email string) (*models.GoogleUser, error) {
	query := `
        SELECT id, email, verified_email, name, given_name, family_name, picture, locale
        FROM users WHERE email = $1`

	row := r.DB.QueryRowContext(ctx, query, email)

	var user models.GoogleUser
	err := row.Scan(&user.ID, &user.Email, &user.VerifiedMail, &user.Name, &user.GivenName, &user.FamilyName, &user.Picture, &user.Locale)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}
