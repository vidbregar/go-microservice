package logger

import (
	"github.com/vidbregar/go-microservice/pkg/config"
	"go.uber.org/zap"
)

func New(config *config.Logger) *zap.Logger {
	if config.Development {
		logger, err := zap.NewDevelopment()
		if err != nil {
			panic("Error initializing dev logger")
		}
		return logger
	}

	logger, err := zap.NewProduction()
	if err != nil {
		panic("Error initializing logger")
	}
	return logger
}
