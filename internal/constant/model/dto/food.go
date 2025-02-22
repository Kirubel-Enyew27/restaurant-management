package dto

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type FoodUpdate struct {
	MealID    uuid.UUID       `json:"meal_id"`
	Name      string          `json:"name,omitempty"`
	Price     decimal.Decimal `json:"price,omitempty"`
	ImgUrl    string          `json:"img_url,omitempty"`
	Available sql.NullBool    `json:"available,omitempty"`
}
