package response

import (
	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(ctx *gin.Context, status int, message string, data interface{}) {
	ctx.JSON(status, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Error(ctx *gin.Context, status int, message string) {
	ctx.JSON(status, Response{
		Success: false,
		Message: message,
	})
}
