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

	order, err := o.orderModule.CreateOrder(ctx, req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.SendSuccessResponse(c, http.StatusCreated, order, nil)
}
func (o *order) GetOrders(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), o.contextTimeout)
	defer cancel()

	orders, err := o.orderModule.GetOrders(ctx)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.SendSuccessResponse(c, http.StatusOK, orders, nil)
}

func (o *order) UpdateOrder(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), o.contextTimeout)
	defer cancel()

	var req dto.CreateOrderRequest

	orderID := c.Param("id")

	if err := c.ShouldBindJSON(&req); err != nil {
		err := errors.ErrBadRequest.Wrap(err, "failed to bind request body")
		o.logger.Info("invalid request body", zap.Error(err))
		_ = c.Error(err)
		return
	}

	updatedOrder, err := o.orderModule.UpdateOrder(ctx, orderID, req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.SendSuccessResponse(c, http.StatusOK, updatedOrder, nil)
}

func (o *order) DeleteOrder(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), o.contextTimeout)
	defer cancel()

	orderID := c.Param("id")

	err := o.orderModule.DeleteOrder(ctx, orderID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.SendSuccessResponse(c, http.StatusOK, err, nil)
}
