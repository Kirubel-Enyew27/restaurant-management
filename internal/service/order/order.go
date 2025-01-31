package order

import (
	"context"
	"restaurant/internal/constant/model/dto"
	"restaurant/internal/service"
	"restaurant/internal/storage"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type clientOrder struct {
	log          *zap.Logger
	storage      storage.Order
	cacheStorage storage.OrderCache
	userStorage  storage.Customer
	foodStorage  storage.Food
}

func InitModule(
	log *zap.Logger,
	storage storage.Order,
	cache storage.OrderCache,
	user storage.Customer,
	food storage.Food,
) service.Order {
	return &clientOrder{
		log:          log,
		storage:      storage,
		cacheStorage: cache,
		userStorage:  user,
		foodStorage:  food,
	}
}

func (o *clientOrder) CreatedOrder(ctx context.Context, order dto.CreateOrderRequest) (dto.OrderResponse, error) {
	user, err := o.userStorage.GetCustomerByID(ctx, order.UserID.UUID)
	if err != nil {
		return dto.OrderResponse{}, err
	}

	var price decimal.Decimal

	for _, item := range order.Item {
		_, err := o.foodStorage.GetFoodByID(ctx, item.MealID.UUID)
		if err != nil {
			return dto.OrderResponse{}, err
		}

		price = price.Add(item.Price.Mul(decimal.NewFromInt32(item.Quantity.Int32)))
	}

	createdOrder, err := o.storage.CreatedOrder(ctx, dto.CreateOrderRequest{
		UserID:      order.UserID,
		OrderStatus: order.OrderStatus,
		TotalPrice:  price,
	})
	if err != nil {
		return dto.OrderResponse{}, err
	}

	for _, orderItem := range order.Item {
		_, err := o.storage.CreateOrderItem(ctx, dto.OrderItem{
			OrderItemID: orderItem.OrderItemID,
			OrderID:     uuid.NullUUID{UUID: createdOrder.OrderID, Valid: true},
			MealID:      orderItem.MealID,
			Quantity:    orderItem.Quantity,
			Price:       orderItem.Price,
		})

		if err != nil {
			return dto.OrderResponse{}, err
		}
	}

	orderResponse := dto.OrderResponse{
		OrderID:     createdOrder.OrderID,
		OrderStatus: createdOrder.OrderStatus,
		TotalPrice:  createdOrder.TotalPrice,
		User:        user,
		OrderItem:   order.Item,
		CreatedAt:   createdOrder.CreatedAt,
		ModifiedAt:  createdOrder.ModifiedAt,
	}

	return orderResponse, nil

}
