package service_auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sso/internal/domain/models"
	internal_errors "sso/internal/errors"
	"sso/internal/utils"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredetials = errors.New("invalid credentials")
)

type AuthService struct {
	log            *slog.Logger
	authRepository AuthRepository
	tokenTTL       time.Duration
}

type AuthRepository interface {
	SaveUser(
		ctx context.Context,
		email string,
		passHash []byte,
	) (userID int64, err error)

	GetUser(
		ctx context.Context,
		email string,
	) (models.User, error)

	GetApp(
		ctx context.Context,
		appID int,
	) (models.App, error)
}

func NewAuthService(logger *slog.Logger, tokenTTL time.Duration, authRepository AuthRepository) *AuthService {
	return &AuthService{
		log:            logger,
		tokenTTL:       tokenTTL,
		authRepository: authRepository,
	}
}

func (s *AuthService) Login(
	ctx context.Context,
	email string,
	password string,
	app_id int32,
) (token string, err error) {
	const op = "service_auth.Login"

	log := s.log.With(
		slog.String("op", op),
		slog.String("username", email),
	)

	log.Info("attemptig to login user")

	user, err := s.authRepository.GetUser(ctx, email)
	if err != nil {
		if errors.Is(err, internal_errors.ErrUserNotFound) {
			s.log.Warn("user not found" + err.Error())

			return "", fmt.Errorf("%s : %w", op, ErrInvalidCredetials)
		}

		s.log.Error("failed to get user" + err.Error())

		return "", fmt.Errorf("%s : %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		s.log.Info("ivalid credentials" + err.Error())

		return "", fmt.Errorf("%s : %w", op, err)
	}

	app, err := s.authRepository.GetApp(ctx, int(app_id))
	if err != nil {
		return "", fmt.Errorf("%s : %w", op, err)
	}

	log.Info("user logger is successfly")

	token, err = utils.NewToken(user, app, s.tokenTTL)
	if err != nil {
		s.log.Error("failed to generate token" + err.Error())

		return "", fmt.Errorf("%s : %w", op, err)
	}

	return token, nil
}

func (s *AuthService) RegisterNewUser(
	ctx context.Context,
	email string,
	password string,
) (userID int64, err error) {
	const op = "service_auth.RegisterNewUser"

	log := s.log.With(
		slog.String("op", op),
		slog.String("username", email),
	)

	log.Info("Register User:" + email)

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("failed to pashed password:%w", err)
	}

	userID, err = s.authRepository.SaveUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, internal_errors.ErrUserExists) {
			log.Warn("user already exists")
			return 0, fmt.Errorf("%s : %w", op, err)
		}
		log.Error("failed to create new user:" + err.Error())
		return 0, fmt.Errorf("%s:%w", op, err)
	}

	log.Info("user is registred")

	return userID, nil

}
