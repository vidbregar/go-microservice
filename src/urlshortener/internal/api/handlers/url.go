package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/vidbregar/go-microservice/internal/db/models"
	"github.com/vidbregar/go-microservice/internal/db/redis/urlshortener"
	"github.com/vidbregar/go-microservice/pkg/shortpath"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"path"
)

const shortPathLen = 7

type UrlHandler interface {
	RedirectUrl(c *gin.Context)
	CreateUrl(c *gin.Context)
}

type handler struct {
	router *gin.RouterGroup
	db     urlshortener.Storage
	gen    shortpath.Generator
	logger *zap.Logger
}

func NewUrlHandler(router *gin.RouterGroup, db urlshortener.Storage, gen shortpath.Generator, logger *zap.Logger) *handler {
	h := handler{
		router: router.Group("/url"),
		db:     db,
		gen:    gen,
		logger: logger,
	}

	h.router.GET("/:shortPath", h.RedirectUrl)
	h.router.POST("/", h.CreateUrl)

	return &h
}

func (h *handler) RedirectUrl(c *gin.Context) {
	shortPath := c.Param("shortPath")

	if !isValidShortPath(shortPath) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid short path"})
		return
	}

	item, err := h.db.Load(c, shortPath)
	if err != nil {
		h.logger.Error(err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	if item != nil && item.Url != "" {
		c.Redirect(http.StatusFound, item.Url)
		return
	}

	c.Status(http.StatusNotFound)
}

func (h *handler) CreateUrl(c *gin.Context) {
	var item models.UrlItem
	if err := c.ShouldBindJSON(&item); err != nil {
		h.logger.Error(err.Error())
		c.Status(http.StatusBadRequest)
		return
	}

	if !isValidUrlItem(&item) {
		c.Status(http.StatusBadRequest)
		return
	}

	var shortPath string
	var err error
	for do := true; do; do = errors.Is(err, urlshortener.ErrShortPathExists) {
		shortPath = h.gen.Generate(shortPathLen)
		err = h.db.Save(c, shortPath, &item)
	}
	if err != nil {
		h.logger.Error(err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	// TODO: Create a model
	shortUrl := path.Join(c.Request.Host, h.router.BasePath(), shortPath)

	c.JSON(http.StatusCreated, gin.H{"shortUrl": shortUrl})
}

func isValidShortPath(shortPath string) bool {
	return len(shortPath) == shortPathLen
}

func isValidUrlItem(item *models.UrlItem) bool {
	if item == nil {
		return false
	}

	_, err := url.ParseRequestURI(item.Url)
	if err != nil {
		return false
	}

	return true
}
