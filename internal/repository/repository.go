package repository_auth

import (
	"context"
	"database/sql"
	"fmt"
	"sso/internal/domain/models"
	internal_errors "sso/internal/errors"
	"strings"

	_ "modernc.org/sqlite"
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(storagePath string) *AuthRepository {
	const op = "repository_auth.NewAuthRepository"

	db, err := sql.Open("sqlite", storagePath)
	if err != nil {
		panic(err)
	}

	return &AuthRepository{
		db: db,
	}
}

func (r *AuthRepository) SaveUser(
	ctx context.Context,
	email string,
	passHash []byte,
) (userID int64, err error) {
	const op = "repository_auth.SaveUser"

	err = r.db.QueryRowContext(ctx, "INSERT INTO users(email, pass_hash) VALUES(?, ?) RETURNING id", email, passHash).Scan(&userID)

	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return 0, internal_errors.ErrUserExists
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return userID, nil
}

func (r *AuthRepository) GetUser(
	ctx context.Context,
	email string,
) (models.User, error) {
	const op = "repository_auth.GetUser"
	row := r.db.QueryRowContext(ctx, "SELECT id, email, pass_hash FROM users WHERE email = ?", email)

	var user models.User

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.PassHash,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, internal_errors.ErrUserNotFound
		}
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil

}

func (r *AuthRepository) GetApp(
	ctx context.Context,
	appID int,
) (models.App, error) {
	const op = "repository_auth.GetApp"
	row := r.db.QueryRowContext(ctx, "SELECT id, name, secret FROM apps WHERE id = ?", appID)

	var app models.App

	err := row.Scan(
		&app.ID,
		&app.Name,
		&app.Secret,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.App{}, internal_errors.ErrAppNotFound
		}

		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}

	return app, nil
}
