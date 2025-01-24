package food

import (
	"restaurant/internal/storage"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

type food struct {
	cache *redis.Client
	log   *zap.Logger
}

func InitFoodCache(cache *redis.Client, OrderExpirationTime time.Duration, log *zap.Logger) storage.FoodCache {
	return &food{
		cache: cache,
		log:   log,
	}
}
