package api

import (
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/vidbregar/go-microservice/internal/api/oapi"
)

type HealthHandler interface {
	GetLivez(ctx echo.Context) error
	GetReadyz(ctx echo.Context) error
}

type healthHandler struct {
	rdb *redis.Client
}

func NewHealthHandler(rdb *redis.Client) HealthHandler {
	return &healthHandler{
		rdb: rdb,
	}
}

func (h *healthHandler) GetLivez(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, oapi.Health{Status: "OK"})
}

func (h *healthHandler) GetReadyz(ctx echo.Context) error {
	_, err := h.rdb.Ping(ctx.Request().Context()).Result()
	if err != nil {
		return ctx.NoContent(http.StatusServiceUnavailable)
	}

	return ctx.JSON(http.StatusOK, oapi.Health{Status: "OK"})
}
