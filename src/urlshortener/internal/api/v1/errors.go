package v1

import (
	"net/http"
	"strconv"

	"github.com/vidbregar/go-microservice/internal/api/oapi"
)

var ErrBadRequest = oapi.BadRequest{
	Code:    strconv.Itoa(http.StatusBadRequest),
	Message: "Request was badly formatted",
}

var ErrNotFound = oapi.NotFound{
	Code:    strconv.Itoa(http.StatusNotFound),
	Message: "Requested resource was not found",
}
