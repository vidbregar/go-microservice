package main

import (
	"github.com/vidbregar/go-microservice/internal/api"
	"github.com/vidbregar/go-microservice/internal/config"
	"github.com/vidbregar/go-microservice/internal/db/redis"
	"github.com/vidbregar/go-microservice/internal/db/redis/urlshortener"
	loggerpkg "github.com/vidbregar/go-microservice/internal/logger"
	"github.com/vidbregar/go-microservice/pkg/shortpath"
	"time"
)

func main() {
	conf := config.Config{
		Server: config.Server{
			":8080",
		},
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

	rdb := redis.New(conf.Redis, logger)
	urlDb := urlshortener.New(rdb)

	gen := shortpath.New(time.Now().UnixNano())

	server := api.New(urlDb, gen, logger)
	server.ListenAndServe(conf.Server.Address)
}
