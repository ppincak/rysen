package errors

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ppincak/rysen/pkg/errors"
)

func ErrorBadRequest(context *gin.Context, err error) {
	var result errors.Error
	if assertion, ok := err.(errors.Error); ok {
		result = assertion
	} else {
		result = errors.NewErrorWithCode("Internal server error", "internal.server.error")
	}
	context.JSON(http.StatusBadRequest, result)
}

func BadRequest(context *gin.Context, err string, key string) {
	context.JSON(http.StatusBadRequest, errors.NewErrorWithCode(err, key))
}

func InternalServerError(context *gin.Context, err string, key string) {
	context.JSON(http.StatusInternalServerError, errors.NewErrorWithCode(err, key))
}
