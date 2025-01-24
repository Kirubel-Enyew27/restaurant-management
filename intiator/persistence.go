package intiator

import (
	"restaurant/internal/constant/model/persistencedb"
	"restaurant/internal/storage"
	"restaurant/internal/storage/persistence/customer"
	"restaurant/internal/storage/persistence/food"
	"restaurant/internal/storage/persistence/order"
	"restaurant/internal/storage/persistence/price"

	"go.uber.org/zap"
)

type Persistence struct {
	food     storage.Food
	order    storage.Order
	customer storage.Customer
	price    storage.Price
}

func InitPersistence(db persistencedb.PersistenceDB, log *zap.Logger) Persistence {
	return Persistence{
		food:     food.Init(db, log),
		order:    order.Init(db, log),
		customer: customer.Init(db, log),
		price:    price.Init(db, log),
	}
}
