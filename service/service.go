package service

import (
	"context"
	"database/sql"
	"errors"
	"restaurant/data"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type ServiceInterface interface {
	Register(context context.Context, user data.User) error
}

type Service struct {
	dbQueries *data.Queries
}

func NewService(dbQueries *data.Queries) *Service {
	return &Service{
		dbQueries: dbQueries,
	}
}

func (s *Service) Register(context context.Context, user data.User) error {
	if err := validation.ValidateStruct(&user,
		validation.Field(&user.Username, validation.Required),
		validation.Field(&user.Email, validation.Required, is.Email),
		validation.Field(&user.Password, validation.Required, validation.Length(4, 0)),
	); err != nil {
		return err
	}

	UsernameTaken, err := s.dbQueries.GetUserByUsername(context, user.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
		} else {
			return err
		}
	}

	if UsernameTaken.Username != "" {
		return errors.New("username already exists")
	}

	EmailRegistered, err := s.dbQueries.GetUserByEmail(context, user.Email)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
		} else {
			return err
		}
	}

	if EmailRegistered.Email != "" {
		return errors.New("email already registered")
	}

	userParams := data.CreateUserParams{
		Username: user.Username,
		Password: user.Password,
		Email:    user.Email,
	}

	_, err = s.dbQueries.CreateUser(context, userParams)

	return err

}
