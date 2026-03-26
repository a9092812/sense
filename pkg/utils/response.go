package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func SuccessResponse(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func ErrorResponse(c *gin.Context, code int, err string) {
	c.JSON(code, APIResponse{
		Success: false,
		Error:   err,
	})
}

func InternalServerError(c *gin.Context, err string) {
	ErrorResponse(c, http.StatusInternalServerError, err)
}

func BadRequestError(c *gin.Context, err string) {
	ErrorResponse(c, http.StatusBadRequest, err)
}
