package api

import (
	"github.com/fvbock/endless"
	"github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/vidbregar/go-microservice/internal/api/handlers"
	"github.com/vidbregar/go-microservice/internal/db/redis/urlshortener"
	"github.com/vidbregar/go-microservice/pkg/shortpath"
	"go.uber.org/zap"
	"time"
)

type Server interface {
	ListenAndServe(addr string)
}

type server struct {
	router *gin.Engine
}

func New(urlDb urlshortener.Storage, gen shortpath.Generator, logger *zap.Logger) *server {
	// TODO: User zap as logger
	r := gin.New()

	r.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(logger, true))

	v1 := r.Group("/v1")
	handlers.NewUrlHandler(v1, urlDb, gen, logger)

	return &server{
		router: r,
	}
}

func (s *server) ListenAndServe(addr string) error {
	return endless.ListenAndServe(addr, s.router)
}
