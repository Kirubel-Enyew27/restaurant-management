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
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	token, err := cstmr.customerModule.Login(ctx, db.User{
		Username: req.Username,
		Password: req.Password,
	})

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
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
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
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
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	userID := c.Param("id")

	updatedUser, err := cstmr.customerModule.UpdateUser(ctx, userID, reqBody)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update user",
		})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"message":      "user updated successfully",
		"updated_user": updatedUser,
	})

}
