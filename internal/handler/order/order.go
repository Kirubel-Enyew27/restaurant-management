package order

import (
	"context"
	"net/http"
	"restaurant/internal/constant/errors"
	"restaurant/internal/constant/model/dto"
	"restaurant/internal/constant/model/response"
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

func (o *order) CreateOrder(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), o.contextTimeout)
	defer cancel()

	var req dto.CreateOrderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		err := errors.ErrBadRequest.Wrap(err, "failed to bind request body")
		o.logger.Info("invalid request body", zap.Error(err))
		_ = c.Error(err)
		return
	}

	order, err := o.orderModule.CreatedOrder(ctx, req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.SendSuccessResponse(c, http.StatusCreated, order, nil)
}
func (o *order) GetOrders(ctx *gin.Context) {}
