package dto

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type OrderItem struct {
	OrderItemID uuid.UUID       `json:"order_item_id,omitempty"`
	OrderID     uuid.UUID       `json:"order_id,omitempty"`
	MealID      uuid.UUID       `json:"meal_id,omitempty"`
	Quantity    int32           `json:"quantity,omitempty"`
	Price       decimal.Decimal `json:"price,omitempty"`
}

type CreateOrderRequest struct {
	OrderID     uuid.UUID       `json:"order_id,omitempty"`
	UserID      uuid.UUID       `json:"user_id,omitempty"`
	OrderStatus string          `json:"order_status,omitempty"`
	TotalPrice  decimal.Decimal `json:"total_price,omitempty"`
	Item        []OrderItem
}

type OrderResponse struct {
	OrderID     uuid.UUID       `json:"order_id,omitempty"`
	OrderStatus string          `json:"order_status,omitempty"`
	TotalPrice  decimal.Decimal `json:"total_price,omitempty"`
	User        User
	OrderItem   []OrderItem
	Meal        Meal
	CreatedAt   sql.NullTime
	ModifiedAt  sql.NullTime
}

type User struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Password string    `json:"password,omitempty"`
}

type Meal struct {
	MealID    uuid.UUID       `json:"meal_id"`
	Name      string          `json:"name"`
	Price     decimal.Decimal `json:"price"`
	Available *bool           `json:"available,omitempty"`
}
