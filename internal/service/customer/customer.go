package customer

import (
	"context"
	"database/sql"
	"restaurant/internal/constant"
	"restaurant/internal/constant/errors"
	"restaurant/internal/constant/model/db"
	"restaurant/internal/constant/model/dto"
	"restaurant/internal/service"
	"restaurant/internal/storage"
	"restaurant/utils"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
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
		c.logger.Error("failed to validate user input", zap.Error(err))
		return db.User{}, errors.ErrInvalidUserInput.Wrap(err, "validation failed")
	}

	existingUser, err := c.storage.GetUserByUsername(ctx, user.Username)
	if err != nil && !strings.Contains(err.Error(), pgx.ErrNoRows.Error()) {
		return db.User{}, err
	} else if existingUser.Username != "" {
		c.logger.Error("username already exists", zap.Error(err))
		return db.User{}, errors.ErrDataAlredyExist.Wrap(err, "username already exists")
	}

	existingUser, err = c.storage.GetUserByEmail(ctx, user.Email)
	if err != nil && !strings.Contains(err.Error(), pgx.ErrNoRows.Error()) {
		return db.User{}, err
	} else if existingUser.Email != "" {
		c.logger.Error("email already exists", zap.Error(err))
		return db.User{}, errors.ErrDataAlredyExist.Wrap(err, "email already exists")
	}

	hashedPassword, err := utils.HashPassword(user.Password, c.logger)
	if err != nil {
		return db.User{}, err
	}
	user.Password = hashedPassword

	return c.storage.Register(ctx, user)
}

func (c *Customer) Login(ctx context.Context, user db.User) (string, error) {
	registeredUser, err := c.storage.GetUserByUsername(ctx, user.Username)
	if err != nil {
		return "", err
	}

	err = utils.VerifyPassword(registeredUser.Password, user.Password, c.logger)
	if err != nil {
		return "", err
	}

	expirationTime := time.Now().Add(constant.TokenExpiration)
	token, err := utils.GenerateJWT(registeredUser, expirationTime, c.logger)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (c *Customer) GetUsers(ctx context.Context) ([]db.User, error) {
	return c.storage.GetCustomers(ctx)
}

func (c *Customer) GetUserByID(ctx context.Context, userID string) (db.User, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		c.logger.Error("failed to parse user id", zap.Error(err))
		return db.User{}, errors.ErrInvalidUserInput.Wrap(err, "invalid user id")
	}

	return c.storage.GetCustomerByID(ctx, userUUID)

}

func (c *Customer) UpdateUser(ctx context.Context, userID string, req dto.UpdateRequest) (db.User, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		c.logger.Error("failed to parse user id", zap.Error(err))
		return db.User{}, errors.ErrInvalidUserInput.Wrap(err, "invalid user id")
	}

	if req.Password != "" {
		hashedPassword, err := utils.HashPassword(req.Password, c.logger)
		if err != nil {
			return db.User{}, err
		}
		req.Password = hashedPassword
	}

	user := db.User{
		UserID:         userUUID,
		Username:       req.Username,
		Password:       req.Password,
		Email:          req.Email,
		ProfilePicture: sql.NullString{String: req.ProfilePicture},
	}

	return c.storage.UpdateCustomer(ctx, user)
}

func (c *Customer) DeleteUser(ctx context.Context, userID string) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		c.logger.Error("failed to parse user id", zap.Error(err))
		return errors.ErrInvalidUserInput.Wrap(err, "invalid user id")
	}

	user, err := c.storage.GetCustomerByID(ctx, userUUID)
	if err != nil {
		return err
	}

	return c.storage.DeleteCustomer(ctx, user.UserID)
}

func (c *Customer) ChangePassword(ctx context.Context, userID string, pass dto.ChangePassword) (db.User, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		c.logger.Error("failed to parse user id", zap.Error(err))
		return db.User{}, errors.ErrInvalidUserInput.Wrap(err, "invalid user id")
	}

	if pass.OldPassword == pass.NewPassword {
		return db.User{}, errors.ErrInvalidUserInput.New("new password cannot be the same as the old password")
	}

	user, err := c.storage.GetCustomerByID(ctx, userUUID)
	if err != nil {
		return db.User{}, err
	}

	err = utils.VerifyPassword(user.Password, pass.OldPassword, c.logger)
	if err != nil {
		return db.User{}, err
	}

	newHashedPassword, err := utils.HashPassword(pass.NewPassword, c.logger)
	if err != nil {
		return db.User{}, err
	}

	updatedUser := db.User{
		UserID:   userUUID,
		Password: newHashedPassword,
	}

	return c.storage.UpdateCustomer(ctx, updatedUser)
}
