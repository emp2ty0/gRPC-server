package repository_auth

import (
	"context"
	"fmt"
	"sso/internal/domain/models"
)

type AuthRepository struct {
}

func NewAuthRepository() *AuthRepository {
	return &AuthRepository{}
}

func (r *AuthRepository) SaveUser(
	ctx context.Context,
	email string,
	passHash []byte,
) (userID int64, err error) {
	return 0, fmt.Errorf("implement me")
}

func (r *AuthRepository) GetUser(
	ctx context.Context,
	email string,
) (models.User, error) {
	return models.User{}, fmt.Errorf("implement me")
}

func (r *AuthRepository) GetApp(
	ctx context.Context,
	appID int,
) (models.App, error) {
	return models.App{}, fmt.Errorf("implement me")
}
