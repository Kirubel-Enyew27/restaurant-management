package order

import (
	"context"
	"restaurant/internal/constant/errors"
	"restaurant/internal/constant/model/dto"
	"restaurant/internal/service"
	"restaurant/internal/storage"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
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
	if err := validation.ValidateStruct(&order,
		validation.Field(&order.UserID, validation.Required),
		validation.Field(&order.Item, validation.Required),
		validation.Field(&order.TotalPrice, validation.Required),
	); err != nil {
		o.log.Error("failed to validate input", zap.Error(err))
		return dto.OrderResponse{}, errors.ErrInvalidUserInput.Wrap(err, "validation failed")
	}

	orderUUID, err := uuid.Parse(order.OrderID.UUID.String())
	if err != nil {
		o.log.Error("failed to parse order id", zap.Error(err))
		return dto.OrderResponse{}, errors.ErrInvalidUserInput.Wrap(err, "invalid order id")
	}

	existingOrder, err := o.storage.GetOrderByID(ctx, orderUUID)
	if err != nil && !strings.Contains(err.Error(), pgx.ErrNoRows.Error()) {
		return dto.OrderResponse{}, err
	} else if existingOrder.OrderStatus.String != "" {
		o.log.Error("order already created", zap.Error(err))
		return dto.OrderResponse{}, errors.ErrDataAlredyExist.Wrap(err, "order already created")
	}

	for _, item := range order.Item {
		orderItemUUID, err := uuid.Parse(item.OrderItemID.UUID.String())
		if err != nil {
			o.log.Error("failed to parse order item id", zap.Error(err))
			return dto.OrderResponse{}, errors.ErrInvalidUserInput.Wrap(err, "invalid order item id")
		}

		existingOrderItem, err := o.storage.GetOrderItemByID(ctx, orderItemUUID)
		if err != nil && !strings.Contains(err.Error(), pgx.ErrNoRows.Error()) {
			return dto.OrderResponse{}, err
		} else if existingOrderItem.OrderID.Valid {
			o.log.Error("order item already created", zap.Error(err))
			return dto.OrderResponse{}, errors.ErrDataAlredyExist.Wrap(err, "order item already created")
		}

	}

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

	var orderItems []dto.OrderItem
	for _, orderItem := range order.Item {
		item, err := o.storage.CreateOrderItem(ctx, dto.OrderItem{
			OrderItemID: orderItem.OrderItemID,
			OrderID:     uuid.NullUUID{UUID: createdOrder.OrderID, Valid: true},
			MealID:      orderItem.MealID,
			Quantity:    orderItem.Quantity,
			Price:       orderItem.Price,
		})

		if err != nil {
			return dto.OrderResponse{}, err
		}

		orderItems = append(orderItems, dto.OrderItem{
			OrderItemID: uuid.NullUUID{UUID: item.OrderItemID, Valid: true},
			OrderID:     item.OrderID,
			MealID:      item.MealID,
			Quantity:    item.Quantity,
			Price:       item.Price,
		})
	}

	orderResponse := dto.OrderResponse{
		OrderID:     createdOrder.OrderID,
		OrderStatus: createdOrder.OrderStatus,
		TotalPrice:  createdOrder.TotalPrice,
		User:        user,
		OrderItem:   orderItems,
		CreatedAt:   createdOrder.CreatedAt,
		ModifiedAt:  createdOrder.ModifiedAt,
	}

	return orderResponse, nil

}
