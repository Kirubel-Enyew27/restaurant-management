package utils

import (
	"fmt"
	"os"
	"restaurant/internal/constant/errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func HashPassword(password string, logger *zap.Logger) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("failed to hash password", zap.Error(err))
		return "", errors.ErrFailedToHash.Wrap(err, "failed to hash password")
	}
	return string(bytes), err
}

func VerifyPassword(hashedPassword, password string, logger *zap.Logger) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		logger.Error("invalid password input", zap.Error(err))
		return errors.ErrInvalidUserInput.Wrap(err, "invalid password")
	}
	return nil
}

func GenerateJWT(username string, expirationTime time.Time, logger *zap.Logger) (string, error) {
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// secretKey := viper.GetString("auth.jwt-key")
	secretKey := os.Getenv("JWT_KEY")
	if secretKey == "" {
		return "", fmt.Errorf("jwtKey is not set in environment variables")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		logger.Error("failed to generate JWT string", zap.Error(err))
		return "", errors.ErrUnableToGet.Wrap(err, "failed to generate JWT")
	}

	return tokenString, nil
}
