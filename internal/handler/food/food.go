package food

import (
	"restaurant/internal/handler"
	"restaurant/internal/service"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type food struct {
	logger         *zap.Logger
	foodModule     service.Food
	orderModule    service.Order
	contextTimeout time.Duration
}

func Init(log *zap.Logger, orderModule service.Order,
	foodModule service.Food,
	contextTimeout time.Duration) handler.Food {
	return &food{
		logger:         log,
		foodModule:     foodModule,
		orderModule:    orderModule,
		contextTimeout: contextTimeout,
	}
}

func (fd *food) AddFood(ctx *gin.Context)  {}
func (fd *food) GetFoods(ctx *gin.Context) {}
