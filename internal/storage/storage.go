package storage

import (
	"context"
	"github.com/google/uuid"
	"restaurant/internal/constant/model/db"
)

type Food interface{}
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
