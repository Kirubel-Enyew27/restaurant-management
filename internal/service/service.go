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
type Order interface {
	CreateOrder(ctx context.Context, order dto.CreateOrderRequest) (dto.OrderResponse, error)
	GetOrders(ctx context.Context) ([]dto.OrderResponse, error)
	UpdateOrder(ctx context.Context, orderID string, order dto.CreateOrderRequest) (dto.OrderResponse, error)
}
type Food interface {
	AddFood(ctx context.Context, meal db.Meal) (db.Meal, error)
	GetFoods(ctx context.Context) ([]db.Meal, error)
	UpdateFood(ctx context.Context, mealID string, req dto.FoodUpdate) (db.Meal, error)
	DeleteFood(ctx context.Context, mealID string) error
}
type Price interface{}
