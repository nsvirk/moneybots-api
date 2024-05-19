package api

import (
	"github.com/labstack/echo/v4"
)

// response represents a success response
type response struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

// errorResponse represents an error response
type errorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// SendResponse sends a success response
func SendResponse(c echo.Context, statusCode int, data interface{}) error {
	response := &response{
		Status: "ok",
		Data:   data,
	}
	return c.JSON(statusCode, response)
}

// SendError sends an error response
func SendError(c echo.Context, statusCode int, message string) error {
	errorResponse := &errorResponse{
		Status:  "error",
		Message: message,
	}
	return c.JSON(statusCode, errorResponse)
}
