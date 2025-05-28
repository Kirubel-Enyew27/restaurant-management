package intiator

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

func InitCache(url string, log *zap.Logger) *redis.Client {
	opts, err := redis.ParseURL(url)
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to parse redis url: %v", err))
	}

	client := redis.NewClient(opts)

	// if _, err := client.Ping(context.Background()).Result(); err != nil {
	// 	log.Fatal(fmt.Sprintf("Failed to ping redis: %v", err))
	// }

	return client
}
