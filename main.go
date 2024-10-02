package main

import (
	"restaurant/config"
	"restaurant/data"
	"restaurant/handler"
	"restaurant/service"

	"github.com/gin-gonic/gin"
)

func init() {
	config.LoadEnv()
	config.ConnectDB()
}

func main() {

	repository := data.New(config.DB)
	service := service.NewService(repository)
	handler := handler.NewHandler(service)

	router := gin.Default()

	router.POST("/register", handler.Register)

	router.Run(":8080")
}
