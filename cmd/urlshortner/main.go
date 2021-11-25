package main

import (
	"github.com/vidbregar/go-microservice/pkg/config"
	loggerpkg "github.com/vidbregar/go-microservice/pkg/logger"
)

func main() {
	conf := config.Config{
		Logger: config.Logger{
			Development: true,
		},
	}

	logger := loggerpkg.New(&conf.Logger)
	_ = logger
}
