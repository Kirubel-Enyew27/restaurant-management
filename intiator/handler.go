package intiator

import (
	"restaurant/internal/handler"
	"restaurant/internal/handler/customer"
	"restaurant/internal/handler/food"
	"restaurant/internal/handler/order"
	"time"

	"go.uber.org/zap"
)

type Handler struct {
	customer handler.Customer
	order    handler.Order
	food     handler.Food
}

func InitHandler(module Module, log *zap.Logger, timeout time.Duration) Handler {

	return Handler{
		food:     food.Init(log, module.Food, module.Food, timeout),
		order:    order.Init(log, module.Order, module.Customer, timeout),
		customer: customer.Init(log, module.Customer, timeout),
	}
}
