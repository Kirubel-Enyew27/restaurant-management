package food

import (
	"context"
	"net/http"
	"restaurant/internal/constant/errors"
	"restaurant/internal/constant/model/db"
	"restaurant/internal/constant/model/dto"
	"restaurant/internal/constant/model/response"
	"restaurant/internal/handler"
	"restaurant/internal/service"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
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

func (fd *food) AddFood(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), fd.contextTimeout)
	defer cancel()

	type Meal struct {
		Name   string          `json:"name" binding:"required"`
		ImgUrl string          `json:"img_url" binding:"required"`
		Price  decimal.Decimal `json:"price" binding:"required"`
	}

	req := Meal{}

	if err := c.ShouldBindJSON(&req); err != nil {
		err := errors.ErrBadRequest.Wrap(err, "failed to bind request body")
		fd.logger.Info("invalid request body", zap.Error(err))
		_ = c.Error(err)
		return
	}

	meal := db.Meal{
		Name:   req.Name,
		Price:  req.Price,
		ImgUrl: req.ImgUrl,
	}

	registeredFood, err := fd.foodModule.AddFood(ctx, meal)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.SendSuccessResponse(c, http.StatusCreated, registeredFood, nil)
}

func (fd *food) GetFoods(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), fd.contextTimeout)
	defer cancel()

	meals, err := fd.foodModule.GetFoods(ctx)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.SendSuccessResponse(c, http.StatusOK, meals, nil)

}

func (fd *food) GetFoodByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), fd.contextTimeout)
	defer cancel()

	foodID := c.Param("id")
	food, err := fd.foodModule.GetFoodByID(ctx, foodID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.SendSuccessResponse(c, http.StatusOK, food, nil)

}

func (fd *food) UpdateFood(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), fd.contextTimeout)
	defer cancel()

	var reqBody dto.FoodUpdate

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		err := errors.ErrBadRequest.Wrap(err, "failed to bind request body")
		fd.logger.Info("invalid request body", zap.Error(err))
		_ = c.Error(err)
		return
	}

	mealID := c.Param("id")

	updatedUser, err := fd.foodModule.UpdateFood(ctx, mealID, reqBody)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.SendSuccessResponse(c, http.StatusOK, updatedUser, nil)

}

func (fd *food) DeleteFood(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), fd.contextTimeout)
	defer cancel()

	mealID := c.Param("id")

	err := fd.foodModule.DeleteFood(ctx, mealID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.SendSuccessResponse(c, http.StatusOK, err, nil)

}

func (fd *food) SearchFood(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), fd.contextTimeout)
	defer cancel()

	query := c.Query("q")
	if query == "" {
		err := errors.ErrBadRequest.New("query param 'q' is required")
		fd.logger.Info("invalid request body", zap.Error(err))
		_ = c.Error(err)
		return
	}

	meals, err := fd.foodModule.SearchFood(ctx, query)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.SendSuccessResponse(c, http.StatusOK, meals, nil)

}
