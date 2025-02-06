package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/imhasandl/grpc-go/sso/internal/domain/models"
	"github.com/imhasandl/grpc-go/sso/internal/lib/jwt"
	"github.com/imhasandl/grpc-go/sso/internal/lib/logger/handlers/sl"
	"github.com/imhasandl/grpc-go/sso/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

//go:generate mockgen -source=auth.go -destination=mocks/auth_mock.go -package=mocks

// Auth is an interface for authentication service.
// It provides methods for user login, registration, and admin check.
// It uses UserSaver, UserProvider, and AppProvider interfaces for data access.
type Auth struct {
	log *slog.Logger
	usrSaver UserSaver
	usrProvider UserProvider
	appProvider AppProvider
	tokenTTL time.Duration
} 

// UserSaver interface defines methods for saving user data.
type UserSaver interface {
	SaveUser(ctx context.Context, email string, passHash []byte) (int64, error)
}

// UserProvider interface defines methods for retrieving user data.
type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

// AppProvider interface defines methods for retrieving application data.
type AppProvider interface {
	App(ctx context.Context, appID int) (models.App, error)
}

// New creates a new Auth service instance.

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

// Login logs in a user given their email, password, and appID.
// It returns a JWT token upon successful login.
func (a *Auth) Login(ctx context.Context, email, password string, appID int) (string, error) {
	const op = "auth.Login"

	log := a.log.With(
		slog.String("op", op),
		slog.String("username", email),
	)

	log.Info("login info")

	user, err := a.usrProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			a.log.Warn("user not found", sl.Err(err))

			return "", fmt.Errorf("%s: %w", op, err)
		}

		log.Error("failed to get user", sl.Err(err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		log.Error("wrong credentianls", sl.Err(err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	app, err := a.appProvider.App(ctx, appID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("User logged in successfully")

	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		log.Error("failed to generate jwt token", sl.Err(err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

// RegisterNewUser registers a new user given their email and password.
// It returns the user ID upon successful registration.
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
		if errors.Is(err, storage.ErrUserExists) {
			log.Warn("user already exists", sl.Err(err))
			
			return 0, fmt.Errorf("%s: %w", op, err)
		}

		log.Error("failed to save user", sl.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return userID, nil
}

// IsAdmin checks if a user with the given ID is an admin.
// It returns a boolean indicating whether the user is an admin and an error.
func (a *Auth) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "auth.IsAdmin"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("is admin info")

	isAdmin, err := a.usrProvider.IsAdmin(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			log.Warn("can't find is user admin", sl.Err(err))
			
			return false, fmt.Errorf("%s: %w", op, err)
		}

		log.Error("can't find is user admin", sl.Err(err))
		return false, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("checked if user is admin", slog.Bool("is_admin", isAdmin))

	return isAdmin, nil
}
