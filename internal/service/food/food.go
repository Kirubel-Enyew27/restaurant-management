package food

import (
	"restaurant/internal/service"
	"restaurant/internal/storage"

	"go.uber.org/zap"
)

type Food struct {
	log          *zap.Logger
	storage      storage.Order
	cacheStorage storage.OrderCache
	priceStorage storage.Price
}

func InitModule(
	log *zap.Logger,
	storage storage.Food,
	cache storage.FoodCache,
	priceStorage storage.Price,
) service.Food {
	return &Food{
		log:          log,
		storage:      storage,
		cacheStorage: cache,
		priceStorage: priceStorage,
	}
}
