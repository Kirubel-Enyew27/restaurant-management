package customer

import (
	"context"
	"fmt"
	"restaurant/internal/constant/model/db"
	"restaurant/internal/constant/model/persistencedb"
	"restaurant/internal/storage"

	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"
)

type customer struct {
	db  persistencedb.PersistenceDB
	log *zap.Logger
}

func Init(db persistencedb.PersistenceDB, log *zap.Logger) storage.Customer {
	return &customer{
		db:  db,
		log: log,
	}
}

func (c *customer) Register(ctx context.Context, user db.User) (db.User, error) {
	arg := db.CreateUserParams{
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
	}

	newUser, err := c.db.Queries.CreateUser(ctx, arg)
	if err != nil {
		return db.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	registeredUser := db.User{
		UserID:    newUser.UserID,
		Username:  newUser.Username,
		Password:  newUser.Password,
		Email:     newUser.Email,
		CreatedAt: newUser.CreatedAt,
	}

	return registeredUser, nil
}

func (c *customer) GetUserByUsername(ctx context.Context, username string) (db.User, error) {
	user, err := c.db.Queries.GetUserByUsername(ctx, username)
	if err != nil && err != pgx.ErrNoRows {
		return db.User{}, fmt.Errorf("failed to fetch user by username: %w", err)
	}

	existingUser := db.User{
		UserID:    user.UserID,
		Username:  user.Username,
		Password:  user.Password,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	return existingUser, nil
}

func (c *customer) GetUserByEmail(ctx context.Context, email string) (db.User, error) {
	user, err := c.db.Queries.GetUserByEmail(ctx, email)
	if err != nil && err != pgx.ErrNoRows {
		return db.User{}, fmt.Errorf("failed to fetch user by email: %w", err)
	}

	existingUser := db.User{
		UserID:    user.UserID,
		Username:  user.Username,
		Password:  user.Password,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	return existingUser, nil
}
