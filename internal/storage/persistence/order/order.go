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
	"github.com/shopspring/decimal"
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

func (o *order) CreateOrder(ctx context.Context, order dto.CreateOrderRequest) (db.Order, error) {
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

func (o *order) GetOrderItemByOrderID(ctx context.Context, orderItemID uuid.NullUUID) ([]db.OrderItem, error) {
	orderItem, err := o.db.Queries.GetOrderItemByOrderID(ctx, orderItemID)
	if err != nil {
		if err == pgx.ErrNoRows {
			o.log.Error("failed to get order item by order id", zap.Error(err))
			return nil, errors.ErrNoRecordFound.Wrap(err, "order item not found")
		}
		o.log.Error("failed to get order item by order id", zap.Error(err))
		return nil, errors.ErrUnableToGet.Wrap(err, "failed to get order item")
	}

	return orderItem, nil
}

func (o *order) GetOrders(ctx context.Context) ([]db.ListOrdersRow, error) {
	orders, err := o.db.Queries.ListOrders(ctx)
	if err != nil {
		if err == pgx.ErrNoRows {
			o.log.Error("failed to get orders", zap.Error(err))
			return nil, errors.ErrUnableToGet.Wrap(err, "orders not found")
		}
		o.log.Error("failed to get orders", zap.Error(err))
		return nil, errors.ErrUnableToGet.Wrap(err, "failed to get orders")
	}

	return orders, nil
}

func (o *order) UpdateOrder(ctx context.Context, order db.Order) (db.Order, error) {
	updateParams := db.UpdateOrderParams{
		OrderID:     order.OrderID,
		OrderStatus: sql.NullString{},
		TotalPrice:  decimal.NullDecimal{},
	}

	if order.OrderStatus.String != "" {
		updateParams.OrderStatus = sql.NullString{String: order.OrderStatus.String, Valid: true}
	}
	if order.TotalPrice.IsZero() == false {
		updateParams.TotalPrice = decimal.NullDecimal{Decimal: order.TotalPrice, Valid: true}
	}

	updatedOrder, err := o.db.Queries.UpdateOrder(ctx, updateParams)
	if err != nil {
		if err == pgx.ErrNoRows {
			o.log.Error("order to be updated does not exist", zap.Error(err))
			return db.Order{}, errors.ErrNoRecordFound.Wrap(err, "order not found")
		}
		o.log.Error("failed to update order", zap.Error(err))
		return db.Order{}, errors.ErrUnableToUpdate.Wrap(err, "failed to update order")
	}

	return updatedOrder, nil
}

func (o *order) UpdateOrderItem(ctx context.Context, orderItem db.OrderItem) (db.OrderItem, error) {
	updateParams := db.UpdateOrderItemParams{
		OrderItemID: orderItem.OrderItemID,
		OrderID:     orderItem.OrderID,
		MealID:      orderItem.MealID,
		Quantity:    orderItem.Quantity,
		Price:       decimal.NullDecimal{},
	}

	if orderItem.Price.IsZero() == false {
		updateParams.Price = decimal.NullDecimal{Decimal: orderItem.Price, Valid: true}
	}

	updatedOrderItem, err := o.db.Queries.UpdateOrderItem(ctx, updateParams)
	if err != nil {
		if err == pgx.ErrNoRows {
			o.log.Error("order item to be updated does not exist", zap.Error(err))
			return db.OrderItem{}, errors.ErrNoRecordFound.Wrap(err, "order item not found")
		}
		o.log.Error("failed to update order item", zap.Error(err))
		return db.OrderItem{}, errors.ErrUnableToUpdate.Wrap(err, "failed to update order item")
	}

	return updatedOrderItem, nil
}

func (o *order) DeleteOrder(ctx context.Context, orderID uuid.UUID) error {
	err := o.db.Queries.DeleteOrder(ctx, orderID)
	if err != nil {
		o.log.Error("error deleting order", zap.Error(err))
		return errors.ErrDBDelError.Wrap(err, "failed to delete order")
	}

	return nil
}
