package storage

import (
	"context"
	"restaurant/internal/constant/model/db"

	"github.com/google/uuid"
)

type Food interface {
	AddFood(ctx context.Context, meal db.Meal) (db.Meal, error)
	GetFoods(ctx context.Context) ([]db.Meal, error)
	GetFoodByName(ctx context.Context, name string) (db.Meal, error)
	UpdateFood(ctx context.Context, meal db.Meal) (db.Meal, error)
}

type Order interface{}
type Customer interface {
	Register(ctx context.Context, user db.User) (db.User, error)
	GetUserByUsername(ctx context.Context, username string) (db.User, error)
	GetUserByEmail(ctx context.Context, email string) (db.User, error)
	GetCustomers(ctx context.Context) ([]db.User, error)
	GetCustomerByID(ctx context.Context, userID uuid.UUID) (db.User, error)
	UpdateCustomer(ctx context.Context, user db.User) (db.User, error)
	DeleteCustomer(ctx context.Context, userID uuid.UUID) error
}
type Price interface{}
type FoodCache interface{}
type OrderCache interface{}
