package order

import (
	"restaurant/internal/service"
	"restaurant/internal/storage"

	"go.uber.org/zap"
)

type clientOrder struct {
	log          *zap.Logger
	storage      storage.Order
	cacheStorage storage.OrderCache
	priceStorage storage.Price
}

func InitModule(
	log *zap.Logger,
	storage storage.Order,
	cache storage.OrderCache,
	priceStorage storage.Price,
) service.Order {
	return &clientOrder{
		log:          log,
		storage:      storage,
		cacheStorage: cache,
		priceStorage: priceStorage,
	}
}
