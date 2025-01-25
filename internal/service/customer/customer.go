package customer

import (
	"context"
	"errors"
	"fmt"
	"restaurant/internal/constant"
	"restaurant/internal/constant/model/db"
	"restaurant/internal/service"
	"restaurant/internal/storage"
	"restaurant/utils"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"go.uber.org/zap"
)

type Customer struct {
	logger  *zap.Logger
	storage storage.Customer
}

func InitModule(
	log *zap.Logger,
	customerStorage storage.Customer,
) service.Customer {
	return &Customer{
		logger:  log,
		storage: customerStorage,
	}

}

func (c *Customer) Register(ctx context.Context, user db.User) (db.User, error) {
	if err := validation.ValidateStruct(&user,
		validation.Field(&user.Username, validation.Required),
		validation.Field(&user.Email, validation.Required, is.Email),
		validation.Field(&user.Password, validation.Required, validation.Length(4, 0)),
	); err != nil {
		return db.User{}, err
	}

	if existingUser, err := c.storage.GetUserByUsername(ctx, user.Username); err != nil {
		return db.User{}, err
	} else if existingUser.Username != "" {
		return db.User{}, errors.New("username already exists")
	}

	if existingUser, err := c.storage.GetUserByEmail(ctx, user.Email); err != nil {
		return db.User{}, err
	} else if existingUser.Email != "" {
		return db.User{}, errors.New("email already registered")
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return db.User{}, err
	}
	user.Password = hashedPassword

	registeredUser, err := c.storage.Register(ctx, user)
	if err != nil {
		return db.User{}, err
	}

	return registeredUser, nil
}

func (c *Customer) Login(ctx context.Context, user db.User) (string, error) {
	registeredUser, err := c.storage.GetUserByUsername(ctx, user.Username)
	if err != nil {
		return "", fmt.Errorf("failed to get user by username: %w", err)
	}

	if registeredUser.Username == "" || !utils.VerifyPassword(registeredUser.Password, user.Password) {
		return "", errors.New("invalid username or password")
	}

	expirationTime := time.Now().Add(constant.TokenExpiration)
	token, err := utils.GenerateJWT(user.UserID.UUID.String(), expirationTime)
	if err != nil {
		return "", fmt.Errorf("failed to generate JWT: %w", err)
	}

	return token, nil
}

func (c *Customer) GetCustomers(ctx context.Context) ([]db.User, error) {
	return c.storage.GetCustomers(ctx)
}
