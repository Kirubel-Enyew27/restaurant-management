package food

import (
	"net/http"
	"restaurant/internal/handler"
	"restaurant/internal/routing"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func InitRoute(group *gin.RouterGroup, h handler.Food,
	log *zap.Logger,
) {
	foodRoutes := []routing.Router{
		{
			Method:      http.MethodPost,
			Path:        "/food/add",
			Handler:     h.AddFood,
			Middlewares: []gin.HandlerFunc{},
		},
		{
			Method:      http.MethodGet,
			Path:        "/foods",
			Handler:     h.GetFoods,
			Middlewares: []gin.HandlerFunc{},
		},
		{
			Method:      http.MethodPatch,
			Path:        "/food/update/:id",
			Handler:     h.AddFood,
			Middlewares: []gin.HandlerFunc{},
		},
		{
			Method:      http.MethodDelete,
			Path:        "/food/delete/:id",
			Handler:     h.GetFoods,
			Middlewares: []gin.HandlerFunc{},
		},
	}

	routing.RegisterRoute(group, foodRoutes, log)
}
