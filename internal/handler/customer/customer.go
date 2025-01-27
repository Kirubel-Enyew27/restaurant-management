package customer

import (
	"context"
	"net/http"
	"restaurant/internal/constant/model/db"
	"restaurant/internal/constant/model/dto"
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
		cstmr.logger.Info("invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body: " + err.Error(),
		})
		return
	}

	newUser, err := cstmr.customerModule.Register(ctx, db.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		cstmr.logger.Info("registration failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "registration failed: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":         "user registered successfully",
		"registered_user": newUser,
	})
}

func (cstmr *customer) Login(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), cstmr.contextTimeout)
	defer cancel()

	var req db.User

	if err := c.ShouldBindJSON(&req); err != nil {
		cstmr.logger.Info("invalid request body", zap.Error(err))
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "invalid request: " + err.Error(),
		})
		return
	}

	token, err := cstmr.customerModule.Login(ctx, db.User{
		Username: req.Username,
		Password: req.Password,
	})

	if err != nil {
		cstmr.logger.Info("unable to login", zap.Error(err))
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "unable to login: " + err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"message": "logged in successffuly",
		"token":   token,
	})
}

func (cstmr *customer) GetCustomers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), cstmr.contextTimeout)
	defer cancel()

	users, err := cstmr.customerModule.GetCustomers(ctx)
	if err != nil {
		cstmr.logger.Info("failed to fetch customers", zap.Error(err))
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"error": "failed to fetch customers: " + err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"message": "users fetched successffuly",
		"users":   users,
	})
}

func (cstmr *customer) UpdateCustomer(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), cstmr.contextTimeout)
	defer cancel()

	var reqBody dto.UpdateRequest

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		cstmr.logger.Info("invalid request body", zap.Error(err))
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body: " + err.Error(),
		})
		return
	}

	userID := c.Param("id")

	updatedUser, err := cstmr.customerModule.UpdateUser(ctx, userID, reqBody)
	if err != nil {
		cstmr.logger.Info("failed to update user", zap.Error(err))
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update user: " + err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"message":      "user updated successfully",
		"updated_user": updatedUser,
	})

}

func (cstmr *customer) DeleteCustomer(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), cstmr.contextTimeout)
	defer cancel()

	userID := c.Param("id")

	err := cstmr.customerModule.DeleteUser(ctx, userID)
	if err != nil {
		cstmr.logger.Info("failed to delete user", zap.Error(err))
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "failed to delete user: " + err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"message": "user deleted successfully",
	})

}
