package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/vidbregar/go-microservice/internal/config"
	"go.uber.org/zap"
)

func New(conf config.Redis, logger *zap.Logger) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:       fmt.Sprintf("%s:%s", conf.Host, conf.Port),
		Username:   conf.Username,
		Password:   conf.Password,
		DB:         0,
		MaxRetries: conf.Retries,
	})

	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		logger.Fatal(fmt.Sprintf("Unable to connect to %s:%s", conf.Host, conf.Port))
	}

	return rdb
}
