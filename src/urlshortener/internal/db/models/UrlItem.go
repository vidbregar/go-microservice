package models

import (
	"github.com/go-openapi/swag"
	"github.com/vidbregar/go-microservice/internal/api"
)

type URL struct {
	Url      string `redis:"url" json:"url"`
	ExpireAt int64  `redis:"expireAt" json:"expireAt"`
}

func ToSwaggerURL(url *URL) api.URL {
	return api.URL{
		ExpireAt: swag.Int64(url.ExpireAt),
		Url:      url.Url,
	}
}

func FromSwaggerURL(url api.URL) URL {
	return URL{
		Url:      url.Url,
		ExpireAt: swag.Int64Value(url.ExpireAt),
	}
}
