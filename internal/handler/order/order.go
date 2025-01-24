package order

import (
	"restaurant/internal/handler"
	"restaurant/internal/service"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type order struct {
	logger         *zap.Logger
	orderModule    service.Order
	customerModule service.Customer
	contextTimeout time.Duration
}

func Init(log *zap.Logger, orderModule service.Order,
	customerModule service.Customer,
	contextTimeout time.Duration) handler.Order {
	return &order{
		logger:         log,
		orderModule:    orderModule,
		customerModule: customerModule,
		contextTimeout: contextTimeout,
	}
}

func (o *order) CreateOrder(ctx *gin.Context) {}
func (o *order) GetOrders(ctx *gin.Context)   {}
