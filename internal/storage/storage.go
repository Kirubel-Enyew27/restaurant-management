package storage

import (
	"context"
	"restaurant/internal/constant/model/db"
	"restaurant/internal/constant/model/dto"

	"github.com/google/uuid"
)

type Food interface {
	AddFood(ctx context.Context, meal db.Meal) (db.Meal, error)
	GetFoods(ctx context.Context) ([]db.Meal, error)
	GetFoodByID(ctx context.Context, mealID uuid.UUID) (db.Meal, error)
	GetFoodByName(ctx context.Context, name string) (db.Meal, error)
	UpdateFood(ctx context.Context, meal db.Meal) (db.Meal, error)
	DeleteFood(ctx context.Context, mealID uuid.UUID) error
}

type Order interface {
	CreateOrder(ctx context.Context, order dto.CreateOrderRequest) (db.Order, error)
	GetOrders(ctx context.Context) ([]db.ListOrdersRow, error)
	CreateOrderItem(ctx context.Context, orderItem dto.OrderItem) (db.OrderItem, error)
	GetOrderByID(ctx context.Context, orderID uuid.UUID) (db.Order, error)
	GetOrderItemByID(ctx context.Context, orderItemID uuid.UUID) (db.OrderItem, error)
	GetOrderItemByOrderID(ctx context.Context, orderItemID uuid.NullUUID) ([]db.OrderItem, error)
	UpdateOrder(ctx context.Context, order db.Order) (db.Order, error)
	UpdateOrderItem(ctx context.Context, orderItem db.OrderItem) (db.OrderItem, error)
}
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
