package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// IndexHandler is the handler for the index route
func IndexHandler(c echo.Context) error {
	return SendResponse(c, http.StatusOK, fmt.Sprintf("%s %s", API_NAME, API_VERSION))
}
