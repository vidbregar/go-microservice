package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/vidbregar/go-microservice/internal/config"
	"go.uber.org/zap"
)

func New(conf config.Redis, logger *zap.Logger) *redis.Client {
	logger.Info("Connecting to Redis...")

	rdb := redis.NewClient(&redis.Options{
		Addr:       fmt.Sprintf("%s:%s", conf.Host, conf.Port),
		Username:   conf.Username,
		Password:   conf.Password,
		DB:         0,
		MaxRetries: conf.Retries,
		OnConnect: func(ctx context.Context, cn *redis.Conn) error {
			logger.Info("Connected to Redis")
			return nil
		},
		MinRetryBackoff: time.Duration(conf.MinRetryBackoff) * time.Millisecond,
		MaxRetryBackoff: time.Duration(conf.MaxRetryBackoff) * time.Millisecond,
	})

	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		logger.Fatal(fmt.Sprintf("Unable to connect to %s:%s", conf.Host, conf.Port))
	}

	return rdb
}
