package api

import (
	"net/http"
	"strconv"
)

var ErrBadRequest = BadRequest{
	Code:    strconv.Itoa(http.StatusBadRequest),
	Message: "Request was badly formatted",
}

var ErrNotFound = NotFound{
	Code:    strconv.Itoa(http.StatusNotFound),
	Message: "Requested resource was not found",
}
