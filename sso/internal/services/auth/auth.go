package auth

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/imhasandl/grpc-go/sso/internal/domain/models"
	"github.com/imhasandl/grpc-go/sso/internal/lib/logger/handlers/sl"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	log *slog.Logger
	usrSaver UserSaver
	usrProvider UserProvider
	appProvider AppProvider
	tokenTTL time.Duration
} 

type UserSaver interface {
	SaveUser(ctx context.Context, email string, passHash []byte) (int64, error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	isAdmin(ctx context.Context, userID int64) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appID int) (models.App, error)
}

func New(
	log *slog.Logger, 
	userSaver UserSaver,
	userProvider UserProvider, 
	appProvider AppProvider, 
	tokenTTL time.Duration,
	) *Auth {
		return &Auth{
			usrSaver: userSaver,
			usrProvider: userProvider,
			log: log,
			appProvider: appProvider,
			tokenTTL: tokenTTL,
		}
}

func (a *Auth) Login(ctx context.Context, email, password string, appID int) (string, error) {
	panic("not implemented")
}

func (a *Auth) RegisterNewUser(ctx context.Context, email, password string) (int64, error) {
	const op = "auth.RegisterNewUser"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("registering info")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to hash the password", sl.Err(err))

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	userID, err := a.usrSaver.SaveUser(ctx, email, hashedPassword)
	if err != nil {
		log.Error("failed to save user", sl.Err(err))

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return userID, nil
}

func (a *Auth) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	panic("not implemented")
}