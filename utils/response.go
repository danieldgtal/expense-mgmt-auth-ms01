package utils

import (
	"github.com/gin-gonic/gin"
)

// APIResponse is a struct for a standardized API response
type APIResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

// SendResponse creates and sends a standardized JSON response
func SendResponse(context *gin.Context, status int, message string, data interface{}, errors interface{}) {
	context.JSON(status, APIResponse{
		Status:  status,
		Message: message,
		Data:    data,
		Errors:  errors,
	})
}
