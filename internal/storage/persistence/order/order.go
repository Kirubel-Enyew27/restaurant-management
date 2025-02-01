package order

import (
	"context"
	"database/sql"
	"restaurant/internal/constant/errors"
	"restaurant/internal/constant/model/db"
	"restaurant/internal/constant/model/dto"
	"restaurant/internal/constant/model/persistencedb"
	"restaurant/internal/storage"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"
)

type order struct {
	db  persistencedb.PersistenceDB
	log *zap.Logger
}

func Init(db persistencedb.PersistenceDB, log *zap.Logger) storage.Order {
	return &order{
		db:  db,
		log: log,
	}
}

func (o *order) CreatedOrder(ctx context.Context, order dto.CreateOrderRequest) (db.Order, error) {
	createOrderParams := db.CreateOrderParams{
		UserID:      order.UserID,
		OrderStatus: sql.NullString{},
		TotalPrice:  order.TotalPrice,
	}

	if order.OrderStatus != "" {
		createOrderParams.OrderStatus = sql.NullString{String: order.OrderStatus, Valid: true}
	}

	createdOrder, err := o.db.Queries.CreateOrder(ctx, createOrderParams)
	if err != nil {
		o.log.Info("failed to create order", zap.Error(err))
		return db.Order{}, errors.ErrUnableTocreate.Wrap(err, "failed to create order")
	}

	createdOrder = db.Order{
		OrderID:     createdOrder.OrderID,
		UserID:      createOrderParams.UserID,
		OrderStatus: createOrderParams.OrderStatus,
		TotalPrice:  createOrderParams.TotalPrice,
		CreatedAt:   createdOrder.CreatedAt,
		ModifiedAt:  createdOrder.ModifiedAt,
	}

	return createdOrder, nil
}

func (o *order) CreateOrderItem(ctx context.Context, orderItem dto.OrderItem) (db.OrderItem, error) {
	CreateOrderItemParams := db.CreateOrderItemParams{
		OrderID:  orderItem.OrderID,
		MealID:   orderItem.MealID,
		Quantity: orderItem.Quantity,
		Price:    orderItem.Price,
	}

	createdOrderItem, err := o.db.Queries.CreateOrderItem(ctx, CreateOrderItemParams)
	if err != nil {
		o.log.Info("failed to create order item", zap.Error(err))
		return db.OrderItem{}, errors.ErrUnableTocreate.Wrap(err, "failed to create order item")
	}

	createdOrderItem = db.OrderItem{
		OrderItemID: createdOrderItem.OrderItemID,
		OrderID:     createdOrderItem.OrderID,
		MealID:      createdOrderItem.MealID,
		Quantity:    createdOrderItem.Quantity,
		Price:       createdOrderItem.Price,
	}

	return createdOrderItem, nil
}

func (o *order) GetOrderByID(ctx context.Context, orderID uuid.UUID) (db.Order, error) {
	order, err := o.db.Queries.GetOrderByID(ctx, orderID)
	if err != nil {
		if err == pgx.ErrNoRows {
			o.log.Error("failed to get order by id", zap.Error(err))
			return db.Order{}, errors.ErrNoRecordFound.Wrap(err, "order not found")
		}
		o.log.Error("failed to get order by id", zap.Error(err))
		return db.Order{}, errors.ErrUnableToGet.Wrap(err, "failed to get order")
	}

	return order, nil
}

func (o *order) GetOrderItemByID(ctx context.Context, orderItemID uuid.UUID) (db.OrderItem, error) {
	orderItem, err := o.db.Queries.GetOrderItemByID(ctx, orderItemID)
	if err != nil {
		if err == pgx.ErrNoRows {
			o.log.Error("failed to get order item by id", zap.Error(err))
			return db.OrderItem{}, errors.ErrNoRecordFound.Wrap(err, "order item not found")
		}
		o.log.Error("failed to get order item by id", zap.Error(err))
		return db.OrderItem{}, errors.ErrUnableToGet.Wrap(err, "failed to get order item")
	}

	return orderItem, nil
}
