package dto

import (
	"database/sql"
	"restaurant/internal/constant/model/db"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type OrderItem struct {
	OrderItemID uuid.NullUUID   `json:"order_item_id,omitempty"`
	OrderID     uuid.NullUUID   `json:"order_id,omitempty"`
	MealID      uuid.NullUUID   `json:"meal_id,omitempty"`
	Quantity    sql.NullInt32   `json:"quantity,omitempty"`
	Price       decimal.Decimal `json:"price,omitempty"`
}

type CreateOrderRequest struct {
	OrderID     uuid.NullUUID   `json:"order_id,omitempty"`
	UserID      uuid.NullUUID   `json:"user_id,omitempty"`
	OrderStatus string          `json:"order_status,omitempty"`
	TotalPrice  decimal.Decimal `json:"total_price,omitempty"`
	Item        []OrderItem
}

type OrderResponse struct {
	OrderID     uuid.UUID       `json:"order_id,omitempty"`
	OrderStatus sql.NullString  `json:"order_status,omitempty"`
	TotalPrice  decimal.Decimal `json:"total_price,omitempty"`
	User        db.User
	OrderItem   []OrderItem
	CreatedAt   sql.NullTime
	ModifiedAt  sql.NullTime
}
