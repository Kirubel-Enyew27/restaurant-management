package customer

import (
	"net/http"
	"restaurant/internal/handler"
	"restaurant/internal/routing"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func InitRoute(group *gin.RouterGroup, h handler.Customer,
	log *zap.Logger,
) {
	customerRoutes := []routing.Router{
		{
			Method:      http.MethodPost,
			Path:        "/customer/register",
			Handler:     h.Register,
			Middlewares: []gin.HandlerFunc{},
		},
		{
			Method:      http.MethodPost,
			Path:        "/customer/login",
			Handler:     h.Login,
			Middlewares: []gin.HandlerFunc{},
		},
		{
			Method:      http.MethodGet,
			Path:        "/customers",
			Handler:     h.GetCustomers,
			Middlewares: []gin.HandlerFunc{},
		},
	}

	routing.RegisterRoute(group, customerRoutes, log)
}
