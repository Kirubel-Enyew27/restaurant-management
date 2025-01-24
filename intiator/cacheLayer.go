package intiator

import (
	"restaurant/internal/storage"
	"restaurant/internal/storage/cache/food"
	"restaurant/internal/storage/cache/order"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

type CacheLayer struct {
	food  storage.FoodCache
	order storage.OrderCache
}

type CacheOptions struct {
	Redis               *redis.Client
	OrderExpirationTime time.Duration
}

func InitCacheLayer(cacheOptions CacheOptions, log *zap.Logger) CacheLayer {
	return CacheLayer{
		food:  food.InitFoodCache(cacheOptions.Redis, cacheOptions.OrderExpirationTime, log),
		order: order.InitOrderCache(cacheOptions.Redis, cacheOptions.OrderExpirationTime, log),
	}
}
