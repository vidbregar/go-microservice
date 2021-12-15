package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vidbregar/go-microservice/internal/api/oapi"
	"github.com/vidbregar/go-microservice/internal/version"
)

type VersionHandler interface {
	GetV1Version(ctx echo.Context) error
}

type versionHandler struct{}

func NewVersionHandler() VersionHandler {
	return &versionHandler{}
}

func (h *versionHandler) GetV1Version(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, oapi.Version{
		Revision: version.Revision,
		Version:  version.GitTag,
	})
}
