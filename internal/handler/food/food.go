package food

import (
	"context"
	"net/http"
	"restaurant/internal/constant/errors"
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

	// Parse form values
	name := c.PostForm("name")
	priceStr := c.PostForm("price")

	// Parse price into decimal.Decimal
	price, err := decimal.NewFromString(priceStr)
	if err != nil {
		err := errors.ErrBadRequest.Wrap(err, "invalid price format")
		_ = c.Error(err)
		return
	}

	// Parse uploaded image
	file, header, err := c.Request.FormFile("food_picture")
	if err != nil {
		err := errors.ErrBadRequest.Wrap(err, "image file is required")
		_ = c.Error(err)
		return
	}
	defer file.Close()

	formData := dto.FormData{
		Name:   name,
		Price:  price,
		File:   file,
		Header: header,
	}

	registeredFood, err := fd.foodModule.AddFood(ctx, formData)
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
