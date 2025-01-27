package customer

import (
	"context"
	"database/sql"
	"fmt"
	"restaurant/internal/constant/model/db"
	"restaurant/internal/constant/model/persistencedb"
	"restaurant/internal/storage"

	"github.com/google/uuid"
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

func (c *customer) GetCustomers(ctx context.Context) ([]db.User, error) {
	users, err := c.db.Queries.ListUsers(ctx)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("no users found: %w", err)
		}
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}

	fetchedUsers := make([]db.User, len(users))

	for i, user := range users {
		fetchedUsers[i] = db.User{
			UserID:   user.UserID,
			Username: user.Username,
			Password: user.Password,
			Email:    user.Email,
		}
	}

	return fetchedUsers, nil
}

func (c *customer) GetCustomerByID(ctx context.Context, userID uuid.UUID) (db.User, error) {
	user, err := c.db.Queries.GetUserByID(ctx, userID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return db.User{}, fmt.Errorf("user not found: %w", err)
		}
		return db.User{}, fmt.Errorf("failed to fetch user: %w", err)
	}

	return user, nil

}

func (c *customer) UpdateCustomer(ctx context.Context, user db.User) (db.User, error) {
	updateParams := db.UpdateUserParams{
		UserID:   user.UserID,
		Username: sql.NullString{}, // Default to an empty nullable string
		Password: sql.NullString{}, // Default to an empty nullable string
		Email:    sql.NullString{}, // Default to an empty nullable string
	}

	// Set values if non-empty
	if user.Username != "" {
		updateParams.Username = sql.NullString{String: user.Username, Valid: true}
	}
	if user.Password != "" {
		updateParams.Password = sql.NullString{String: user.Password, Valid: true}
	}
	if user.Email != "" {
		updateParams.Email = sql.NullString{String: user.Email, Valid: true}
	}

	updatedUser, err := c.db.Queries.UpdateUser(ctx, updateParams)
	if err != nil {
		return db.User{}, fmt.Errorf("error updating user: %w", err)
	}

	return updatedUser, nil
}

func (c *customer) DeleteCustomer(ctx context.Context, userID uuid.UUID) error {
	err := c.db.Queries.DeleteUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}

	return nil

}
