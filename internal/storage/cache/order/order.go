package order

import (
	"restaurant/internal/storage"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

type order struct {
	cache *redis.Client
	log   *zap.Logger
}

func InitOrderCache(cache *redis.Client, OrderExpirationTime time.Duration, log *zap.Logger) storage.OrderCache {
	return &order{
		cache: cache,
		log:   log,
	}
}
