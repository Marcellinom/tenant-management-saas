package errors

import (
	"errors"
	"github.com/Marcellinom/tenant-management-saas/provider"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type DefaultErrorHandler struct{}

func DefaultHandler() DefaultErrorHandler {
	return DefaultErrorHandler{}
}

type H provider.H

func (e DefaultErrorHandler) Handle(debugMode bool) func(ctx *provider.Context) {
	return func(ctx *provider.Context) {
		ctx.Next()
		err := ctx.Errors.Last()

		if err == nil {
			return
		}
		requestId := ""
		reqIdInterface, exists := ctx.Get("request_id")
		if exists {
			if reqId, ok := reqIdInterface.(string); ok {
				requestId = reqId
			}
		}

		data := gin.H{
			"request_id": requestId,
		}

		var badRequestError BadRequestError
		var unauthorizedError UnauthorizedError
		var invariantError InvariantError
		var expectationFailed ExpectationFailedError

		switch {
		case errors.As(err, &expectationFailed):
			log.Printf("Request ID: %s; Status: %d; Error: %s\n", requestId, http.StatusExpectationFailed, err.Error())
			ctx.JSON(
				http.StatusExpectationFailed,
				H{
					"code":    expectationFailed.Code(),
					"message": expectationFailed.Message(),
					"data":    data,
				},
			)
		case errors.As(err, &badRequestError):
			log.Printf("Request ID: %s; Status: 400; Error: %s\n", requestId, err.Error())
			for key, val := range badRequestError.GetData() {
				data[key] = val
			}
			ctx.JSON(
				http.StatusBadRequest,
				H{
					"code":    badRequestError.Code(),
					"message": badRequestError.Message(),
					"data":    data,
				},
			)
		case errors.As(err, &unauthorizedError):
			log.Printf("Request ID: %s; Status: 401; Error: %s\n", requestId, err.Error())
			ctx.JSON(
				http.StatusUnauthorized,
				H{
					"code":    unauthorizedError.Code(),
					"message": unauthorizedError.Message(),
					"data":    data,
				},
			)
		case errors.As(err, &invariantError):
			log.Printf("Request ID: %s; Status: 500; Error: %s\n", requestId, err.Error())
			ctx.JSON(
				http.StatusInternalServerError,
				H{
					"code":    invariantError.Code(),
					"message": invariantError.Message(),
					"data":    invariantError.Data,
				},
			)
		default:
			log.Printf("Request ID: %s; Status: 500; Error: %s\n", requestId, err.Error())
			if debugMode {
				data["error"] = err.Error()
			}
			ctx.JSON(
				http.StatusInternalServerError,
				H{
					"code": http.StatusInternalServerError,
					"data": data,
				},
			)
		}

		ctx.Abort()
	}
}
