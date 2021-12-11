package api

import (
	"github.com/vidbregar/go-microservice/internal/api/oapi"
	v1 "github.com/vidbregar/go-microservice/internal/api/v1"
)

type server struct {
	v1.UrlHandler
	v1.VersionHandler
}

func NewServer(urlHandlerV1 v1.UrlHandler, versionHandlerV1 v1.VersionHandler) oapi.ServerInterface {
	return &server{
		urlHandlerV1,
		versionHandlerV1,
	}
}
