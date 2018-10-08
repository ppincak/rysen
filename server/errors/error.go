package errors

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ppincak/rysen/pkg/errors"
)

type ApiError struct {
	Message string `json:"message"`
	Key     string `json:"key"`
}

func NewApiError(message string, key string) *ApiError {
	return &ApiError{
		Message: message,
		Key:     key,
	}
}

func ErrorBadRequest(context *gin.Context, err error) {
	var result *ApiError
	if assertion, ok := err.(errors.Error); ok {
		result = NewApiError(assertion.Message, assertion.Code)
	} else {
		result = NewApiError("Internal server error", "internal.server.error")
	}
	context.JSON(http.StatusBadRequest, result)
}

func BadRequest(context *gin.Context, err string, key string) {
	context.JSON(http.StatusBadRequest, NewApiError(err, key))
}

func InternalServerError(context *gin.Context, err string, key string) {
	context.JSON(http.StatusInternalServerError, NewApiError(err, key))
}
