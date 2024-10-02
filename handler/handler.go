package handler

import (
	"context"
	"net/http"
	"restaurant/data"
	"restaurant/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service service.ServiceInterface
}

func NewHandler(service service.ServiceInterface) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(c *gin.Context) {
	var user struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := h.service.Register(context.Background(), data.User{
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user registered successfully",
	})
}
