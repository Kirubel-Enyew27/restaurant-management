package customer

import (
	"context"
	"net/http"
	"restaurant/internal/constant/errors"
	"restaurant/internal/constant/model/db"
	"restaurant/internal/constant/model/dto"
	"restaurant/internal/constant/model/response"
	"restaurant/internal/handler"
	"restaurant/internal/service"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type customer struct {
	logger         *zap.Logger
	customerModule service.Customer
	contextTimeout time.Duration
}

func Init(log *zap.Logger, customerModule service.Customer,
	contextTimeout time.Duration) handler.Customer {
	return &customer{
		logger:         log,
		customerModule: customerModule,
		contextTimeout: contextTimeout,
	}
}

func (cstmr *customer) Register(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), cstmr.contextTimeout)
	defer cancel()

	var req db.User

	if err := c.ShouldBindJSON(&req); err != nil {
		err := errors.ErrBadRequest.Wrap(err, "failed to bind request body")
		cstmr.logger.Info("invalid request body", zap.Error(err))
		_ = c.Error(err)
		return
	}

	newUser, err := cstmr.customerModule.Register(ctx, db.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.SendSuccessResponse(c, http.StatusCreated, newUser, nil)
}

func (cstmr *customer) Login(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), cstmr.contextTimeout)
	defer cancel()

	var req db.User

	if err := c.ShouldBindJSON(&req); err != nil {
		err := errors.ErrBadRequest.Wrap(err, "failed to bind request body")
		cstmr.logger.Info("invalid request body", zap.Error(err))
		_ = c.Error(err)
		return
	}

	token, err := cstmr.customerModule.Login(ctx, db.User{
		Username: req.Username,
		Password: req.Password,
	})

	if err != nil {
		_ = c.Error(err)
		return
	}

	response.SendSuccessResponse(c, http.StatusOK, token, nil)

}

func (cstmr *customer) GetCustomers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), cstmr.contextTimeout)
	defer cancel()

	users, err := cstmr.customerModule.GetUsers(ctx)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.SendSuccessResponse(c, http.StatusOK, users, nil)

}

func (cstmr *customer) UpdateCustomer(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), cstmr.contextTimeout)
	defer cancel()

	var reqBody dto.UpdateRequest

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		err := errors.ErrBadRequest.Wrap(err, "failed to bind request body")
		cstmr.logger.Info("invalid request body", zap.Error(err))
		_ = c.Error(err)
		return
	}

	userID := c.Param("id")

	updatedUser, err := cstmr.customerModule.UpdateUser(ctx, userID, reqBody)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.SendSuccessResponse(c, http.StatusOK, updatedUser, nil)

}

func (cstmr *customer) DeleteCustomer(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), cstmr.contextTimeout)
	defer cancel()

	userID := c.Param("id")

	err := cstmr.customerModule.DeleteUser(ctx, userID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.SendSuccessResponse(c, http.StatusOK, err, nil)

}
