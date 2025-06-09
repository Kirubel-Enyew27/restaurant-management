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
	GetUserByID(ctx context.Context, userID string) (db.User, error)
	UpdateUser(ctx context.Context, param string, req dto.UpdateRequest) (db.User, error)
	DeleteUser(ctx context.Context, userID string) error
	ChangePassword(ctx context.Context, userID string, pass dto.ChangePassword) (db.User, error)
	SearchUser(ctx context.Context, query string) ([]db.User, error)
}
type Order interface {
	CreateOrder(ctx context.Context, order dto.CreateOrderRequest) (dto.OrderResponse, error)
	GetOrders(ctx context.Context) ([]dto.OrderResponse, error)
	UpdateOrder(ctx context.Context, orderID string, order dto.CreateOrderRequest) (dto.OrderResponse, error)
	DeleteOrder(ctx context.Context, orderID string) error
	SearchOrder(ctx context.Context, query string) ([]dto.OrderResponse, error)
}
type Food interface {
	AddFood(ctx context.Context, meal db.Meal) (db.Meal, error)
	GetFoods(ctx context.Context) ([]db.Meal, error)
	GetFoodByID(ctx context.Context, foodID string) (db.Meal, error)
	UpdateFood(ctx context.Context, mealID string, req dto.FoodUpdate) (db.Meal, error)
	DeleteFood(ctx context.Context, mealID string) error
	SearchFood(ctx context.Context, query string) ([]db.Meal, error)
}
type Price interface{}
