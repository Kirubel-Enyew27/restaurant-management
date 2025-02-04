package handler

import (
	"github.com/gin-gonic/gin"
)

type Customer interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	GetCustomers(c *gin.Context)
	UpdateCustomer(c *gin.Context)
	DeleteCustomer(c *gin.Context)
}
type Order interface {
	CreateOrder(c *gin.Context)
	GetOrders(c *gin.Context)
	UpdateOrder(c *gin.Context)
}
type Food interface {
	AddFood(c *gin.Context)
	GetFoods(c *gin.Context)
	UpdateFood(c *gin.Context)
	DeleteFood(c *gin.Context)
}
