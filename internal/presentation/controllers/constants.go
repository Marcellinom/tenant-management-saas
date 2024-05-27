package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var statusCode = map[string]int{
	"success": 1,
}

type DefaultResponse struct {
	Code    int         `json:"code" example:"123"`
	Message string      `json:"message" example:"success"`
	Data    interface{} `json:"data,omitempty" swaggertype:"object"`
}

func SuccessWithData(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusOK, gin.H{
		"code":    statusCode["success"],
		"message": "success",
		"data":    data,
	})
}

func Success(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"code":    statusCode["success"],
		"message": "success",
		"data":    nil,
	})
}
