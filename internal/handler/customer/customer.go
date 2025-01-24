package customer

import (
	"context"
	"net/http"
	"restaurant/internal/constant/model/db"
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
	var req db.User

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body: " + err.Error(),
		})
		return
	}

	newUser, err := cstmr.customerModule.Register(context.Background(), db.User{
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
	var req db.User

	if err := c.ShouldBindJSON(&req); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	token, err := cstmr.customerModule.Login(context.Background(), db.User{
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

func (cstmr *customer) GetCustomers(ctx *gin.Context) {}
