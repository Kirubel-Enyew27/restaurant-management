package order

import (
	"net/http"
	"restaurant/internal/handler"
	"restaurant/internal/routing"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func InitRoute(group *gin.RouterGroup, h handler.Order,
	log *zap.Logger,
) {
	orderRoutes := []routing.Router{
		{
			Method:      http.MethodPost,
			Path:        "/order",
			Handler:     h.CreateOrder,
			Middlewares: []gin.HandlerFunc{},
		},
		{
			Method:      http.MethodGet,
			Path:        "/orders",
			Handler:     h.GetOrders,
			Middlewares: []gin.HandlerFunc{},
		},
	}

	routing.RegisterRoute(group, orderRoutes, log)
}
