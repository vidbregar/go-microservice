package handlers

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/vidbregar/go-microservice/internal/api"
	"github.com/vidbregar/go-microservice/internal/db/models"
	"github.com/vidbregar/go-microservice/internal/db/redis/urlshortener"
	"github.com/vidbregar/go-microservice/pkg/shortpath"
	"go.uber.org/zap"
	"net/http"
	"path"
)

const shortPathLen = 7

type handler struct {
	db     urlshortener.Storage
	gen    shortpath.Generator
	logger *zap.Logger
}

func NewUrlHandler(db urlshortener.Storage, gen shortpath.Generator, logger *zap.Logger) api.ServerInterface {
	return &handler{
		db:     db,
		gen:    gen,
		logger: logger,
	}
}

func (h handler) PostUrl(ctx echo.Context) error {
	var urlSwag api.URL
	if err := ctx.Bind(&urlSwag); err != nil {
		h.logger.Error(err.Error())
		return ctx.JSON(http.StatusBadRequest, api.ErrBadRequest)
	}

	var shortPath string
	var err error
	url := models.FromSwaggerURL(urlSwag)
	for do := true; do; do = errors.Is(err, urlshortener.ErrShortPathExists) {
		shortPath = h.gen.Generate(shortPathLen)
		err = h.db.Save(ctx.Request().Context(), shortPath, &url)
	}
	if err != nil {
		h.logger.Error(err.Error())
		return err
	}

	shortened := api.ShortenedURL{
		ShortUrl: path.Join(ctx.Request().Host, ctx.Path(), shortPath),
	}

	return ctx.JSON(http.StatusCreated, shortened)
}

func (h handler) GetUrlShortened(ctx echo.Context, shortened string) error {
	url, err := h.db.Load(ctx.Request().Context(), shortened)
	if err != nil {
		h.logger.Error(err.Error())
		return err
	}

	if url != nil && url.Url != "" {
		return ctx.Redirect(http.StatusFound, url.Url)
	}

	return ctx.JSON(http.StatusNotFound, api.ErrNotFound)
}
