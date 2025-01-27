package handler

import (
	"github.com/gin-gonic/gin"
)

type Customer interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
	GetCustomers(ctx *gin.Context)
	UpdateCustomer(c *gin.Context)
	DeleteCustomer(c *gin.Context)
}
type Order interface {
	CreateOrder(ctx *gin.Context)
	GetOrders(ctx *gin.Context)
}
type Food interface {
	AddFood(ctx *gin.Context)
	GetFoods(ctx *gin.Context)
}
