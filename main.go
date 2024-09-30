package main

import (
	"restaurant/config"

	"github.com/gin-gonic/gin"
)

func init() {
	config.LoadEnv()
	config.ConnectDB()
}

func main() {

	router := gin.Default()

	router.Run(":8080")
}
