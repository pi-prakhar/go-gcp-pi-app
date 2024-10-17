package repository

import (
	"context"
	"database/sql"
	"log"

	"github.com/pi-prakhar/go-gcp-pi-app/internal/auth/models"
	errors "github.com/pi-prakhar/go-gcp-pi-app/pkg/error"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user models.GoogleUser) error
	GetUserByEmail(ctx context.Context, email string) (*models.GoogleUser, error)
	GetUsers(ctx context.Context) ([]*models.GoogleUser, error)
	DeleteUserByEmail(ctx context.Context, email *string) error
	UpdateUserByEmail(ctx context.Context, user *models.GoogleUser) error
}

type GCPPostgresqlRepository struct {
	DB *sql.DB
}

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

func (r *GCPPostgresqlRepository) GetUserByEmail(ctx context.Context, email string) (*models.GoogleUser, error) {
	query := `
        SELECT id, email, verified_email, name, given_name, family_name, picture, locale
        FROM users WHERE email = $1`

	row := r.DB.QueryRowContext(ctx, query, email)

	var user models.GoogleUser
	err := row.Scan(&user.ID, &user.Email, &user.VerifiedMail, &user.Name, &user.GivenName, &user.FamilyName, &user.Picture, &user.Locale)
	if err == sql.ErrNoRows {
		return nil, errors.ErrUserNotFound
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *GCPPostgresqlRepository) GetUsers(ctx context.Context) ([]*models.GoogleUser, error) {
	var err error
	var users []*models.GoogleUser

	query := `SELECT id, email, verified_email, name, given_name, family_name, picture, locale
        FROM users`

	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.GoogleUser
		err := rows.Scan(&user.ID, &user.Email, &user.VerifiedMail, &user.Name, &user.GivenName, &user.FamilyName, &user.Picture, &user.Locale)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *GCPPostgresqlRepository) UpdateUserByEmail(ctx context.Context, user *models.GoogleUser) error {
	query := `
			UPDATE users
			SET verified_email=$2,
			    name=$3,
			    given_name=$4,
			    family_name=$5,
			    picture=$6,
			    locale=$7
			WHERE email = $1
		`
	res, err := r.DB.ExecContext(ctx, query, user.Email, user.VerifiedMail, user.Name, user.GivenName, user.FamilyName, user.Picture, user.Locale)
	log.Println("res : ", res)
	if err != nil {
		return err
	}
	return nil
}

func (r *GCPPostgresqlRepository) DeleteUserByEmail(ctx context.Context, email *string) error {
	query := `
        DELETE FROM users 
		WHERE email = $1
		`
	res, err := r.DB.ExecContext(ctx, query, *email)
	log.Println("res : ", res)
	if err != nil {
		return nil
	}

	if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
		log.Println("Row Affected : ", rowsAffected)
		return errors.ErrUserNotFound
	} else {
		log.Println("Row Affected : ", rowsAffected)
	}

	return nil
}
