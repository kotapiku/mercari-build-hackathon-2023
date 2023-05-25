package service

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kotapiku/mercari-build-hackathon-2023/backend/db"
	"github.com/kotapiku/mercari-build-hackathon-2023/backend/domain"
	"golang.org/x/crypto/bcrypt"
)

type LoginService struct {
	UserRepo db.UserRepository
}

type JwtCustomClaims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

var (
	ErrMismatchPassword = errors.New("mismatch password")
)

func NewLoginService(sqlDB *sql.DB) LoginService {
	return LoginService{UserRepo: db.NewUserRepository(sqlDB)}
}

func (l LoginService) LoginByID(ctx context.Context, userID int64, password string) (*domain.User, string, error) {
	user, err := l.UserRepo.GetUser(ctx, userID)
	if err != nil {
		return nil, "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return nil, "", ErrMismatchPassword
		}
		return nil, "", err
	}

	// Set custom claims
	claims := &JwtCustomClaims{
		user.ID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}
	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Generate encoded token and send it as response.
	encodedToken, err := token.SignedString([]byte(GetSecret()))
	if err != nil {
		return nil, "", err
	}

	return &user, encodedToken, nil
}

func (l LoginService) LoginByName(ctx context.Context, userName string, password string) (*domain.User, string, error) {
	user, err := l.UserRepo.GetUserByName(ctx, userName)
	if err != nil {
		return nil, "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return nil, "", ErrMismatchPassword
		}
		return nil, "", err
	}

	// Set custom claims
	claims := &JwtCustomClaims{
		user.ID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}
	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Generate encoded token and send it as response.
	encodedToken, err := token.SignedString([]byte(GetSecret()))
	if err != nil {
		return nil, "", err
	}

	return &user, encodedToken, nil
}

func GetSecret() string {
	if secret := os.Getenv("SECRET"); secret != "" {
		return secret
	}
	return "secret-key"
}
