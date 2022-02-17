package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/brpaz/echozap"
	"github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/labstack/echo/v4"
	"github.com/vidbregar/go-microservice/internal/api"
	"github.com/vidbregar/go-microservice/internal/api/oapi"
	"github.com/vidbregar/go-microservice/internal/api/v1"
	"github.com/vidbregar/go-microservice/internal/config"
	"github.com/vidbregar/go-microservice/internal/db/redis"
	"github.com/vidbregar/go-microservice/internal/db/redis/urlshortener"
	loggerpkg "github.com/vidbregar/go-microservice/pkg/logger"
	"github.com/vidbregar/go-microservice/pkg/shortpath"
)

func main() {
	configFile := flag.String("config", "config.yaml", "path to config yaml file")
	secretsDir := flag.String("secrets", "/secrets/", "path to secrets directory")
	flag.Parse()
	conf, err := config.New(*configFile, *secretsDir)
	if err != nil {
		panic(err)
	}

	logger := loggerpkg.New(&conf.Logger)
	defer logger.Sync()

	rdb := redis.New(conf.Redis, logger)
	urlDb := urlshortener.New(rdb)

	gen := shortpath.New(time.Now().UnixNano())

	swagger, err := oapi.GetSwagger()
	if err != nil {
		logger.Fatal(fmt.Sprintf("Error loading swagger spec\n: %s", err))
	}

	e := echo.New()
	e.Use(echozap.ZapLogger(logger))
	e.Use(middleware.OapiRequestValidator(swagger))

	oapi.RegisterHandlers(
		e,
		api.NewServer(
			api.NewHealthHandler(rdb),
			v1.NewUrlHandler(urlDb, gen, logger),
			v1.NewVersionHandler(),
		),
	)

	go func() {
		if err := e.Start(conf.Server.Address); err != nil && err != http.ErrServerClosed {
			logger.Fatal("shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	s := <-quit

	if conf.Server.DelaySigterm > 0 && s == syscall.SIGTERM {
		logger.Info("Delaying SIGTERM")
		time.Sleep(time.Duration(conf.Server.DelaySigterm) * time.Second)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		logger.Fatal(err.Error())
	}
}
