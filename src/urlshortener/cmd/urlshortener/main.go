package main

import (
	"github.com/vidbregar/go-microservice/pkg/config"
	loggerpkg "github.com/vidbregar/go-microservice/pkg/logger"
	"github.com/vidbregar/go-microservice/pkg/storage/redis"
)

func main() {
	conf := config.Config{
		Logger: config.Logger{
			Development: true,
		},
		Redis: config.Redis{
			Host:    "127.0.0.1",
			Port:    "6379",
			Retries: 5,
		},
	}

	logger := loggerpkg.New(&conf.Logger)
	defer logger.Sync()

	client := redis.New(conf.Redis, logger)
	_ = client
}
