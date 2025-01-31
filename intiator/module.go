package intiator

import (
	"restaurant/internal/service"
	"restaurant/internal/service/customer"
	"restaurant/internal/service/food"
	"restaurant/internal/service/order"

	"go.uber.org/zap"
)

type Module struct {
	Customer service.Customer
	Order    service.Order
	Food     service.Food
}

func InitModule(persistence Persistence, cache CacheLayer,
	log *zap.Logger) Module {

	return Module{
		Customer: customer.InitModule(log,
			persistence.customer,
		),
		Order: order.InitModule(
			log, persistence.order, cache.order, persistence.customer, persistence.food),

		Food: food.InitModule(log, persistence.food, cache.food, persistence.price),
	}
}
