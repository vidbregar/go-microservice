package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/brpaz/echozap"
	"github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/fvbock/endless"
	"github.com/labstack/echo/v4"
	"github.com/vidbregar/go-microservice/internal/api"
	"github.com/vidbregar/go-microservice/internal/api/oapi"
	"github.com/vidbregar/go-microservice/internal/api/v1"
	"github.com/vidbregar/go-microservice/internal/config"
	"github.com/vidbregar/go-microservice/internal/db/redis"
	"github.com/vidbregar/go-microservice/internal/db/redis/urlshortener"
	loggerpkg "github.com/vidbregar/go-microservice/internal/logger"
	"github.com/vidbregar/go-microservice/pkg/shortpath"
)

func main() {
	configFile := flag.String("config", "config.yaml", "path to config yaml file")
	flag.Parse()
	conf, err := config.New(*configFile)
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

	err = endless.ListenAndServe(conf.Server.Address, e)
	if err != nil {
		logger.Error(err.Error())
	}
	logger.Info(fmt.Sprintf("Stopped listening on %s", conf.Server.Address))
}
