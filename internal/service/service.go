package service

import (
	"context"
	"restaurant/internal/constant/model/db"
	"restaurant/internal/constant/model/dto"
)

type Customer interface {
	Register(ctx context.Context, user db.User) (db.User, error)
	Login(ctx context.Context, user db.User) (string, error)
	GetUsers(ctx context.Context) ([]db.User, error)
	UpdateUser(ctx context.Context, param string, req dto.UpdateRequest) (db.User, error)
	DeleteUser(ctx context.Context, userID string) error
}
type Order interface{}
type Food interface {
	AddFood(ctx context.Context, meal db.Meal) (db.Meal, error)
	GetFoods(ctx context.Context) ([]db.Meal, error)
}
type Price interface{}
