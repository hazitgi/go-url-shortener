package utils

import (
	"github.com/gin-gonic/gin"
)

type HTTPSuccess struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Status  int         `json:"-"`
}

type HTTPError struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Error   interface{} `json:"errors"`
	Status  int         `json:"-"`
}

func (h *HTTPError) InValidResponse(ctx *gin.Context) {
	ctx.JSON(h.Status, gin.H{
		"success": h.Success,
		"message": h.Message,
		"errors":  h.Error,
	})
}

func (h *HTTPSuccess) SuccessResponse(ctx *gin.Context) {
	ctx.JSON(h.Status, gin.H{
		"success": h.Success,
		"message": h.Message,
		"data":    h.Data,
	})
}
