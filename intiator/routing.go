package intiator

import (
	"restaurant/internal/routing/customer"
	"restaurant/internal/routing/food"
	"restaurant/internal/routing/order"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func InitRouter(
	group *gin.RouterGroup,
	handler Handler,
	module Module,
	log *zap.Logger,
) {

	customer.InitRoute(group, handler.customer, log)
	food.InitRoute(group, handler.food, log)
	order.InitRoute(group, handler.order, log)

}
