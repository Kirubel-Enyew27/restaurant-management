package service

import (
	"context"
	"restaurant/internal/constant/model/db"
)

type Customer interface {
	Register(ctx context.Context, user db.User) (db.User, error)
	Login(ctx context.Context, user db.User) (string, error)
}
type Order interface{}
type Food interface{}
type Price interface{}
