package storage

import (
	"context"
	"restaurant/internal/constant/model/db"
)

type Food interface{}
type Order interface{}
type Customer interface {
	Register(ctx context.Context, user db.User) (db.User, error)
	GetUserByUsername(ctx context.Context, username string) (db.User, error)
	GetUserByEmail(ctx context.Context, email string) (db.User, error)
	GetCustomers(ctx context.Context) ([]db.User, error)
}
type Price interface{}
type FoodCache interface{}
type OrderCache interface{}
