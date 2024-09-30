// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package data

import (
	"database/sql"
)

type Drink struct {
	DrinkID   int32
	Name      string
	Price     string
	Available sql.NullBool
	CreatedAt sql.NullTime
}

type Meal struct {
	MealID    int32
	Name      string
	Price     string
	Available sql.NullBool
	CreatedAt sql.NullTime
}

type Order struct {
	OrderID    int32
	UserID     sql.NullInt32
	TotalPrice string
	Status     string
	OrderDate  sql.NullTime
}

type OrderDrink struct {
	OrderID  int32
	DrinkID  int32
	Quantity int32
}

type OrderMeal struct {
	OrderID  int32
	MealID   int32
	Quantity int32
}

type User struct {
	UserID    int32
	Username  string
	Email     string
	Password  string
	CreatedAt sql.NullTime
}