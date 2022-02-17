package logger

import (
	"github.com/vidbregar/go-microservice/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(config *config.Logger) *zap.Logger {
	if config.Development {
		zapCon := zap.NewDevelopmentConfig()
		zapCon.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		logger, err := zapCon.Build()
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
