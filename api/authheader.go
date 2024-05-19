package api

import (
	"strings"

	"github.com/labstack/echo/v4"
)

// GetAuthHeaderValues returns the userId and enctoken from the authorization header
func GetAuthHeaderValues(c echo.Context) (userId, enctoken string) {

	// Get the Authorization header
	authHeader := c.Request().Header.Get(echo.HeaderAuthorization)

	// Check if the Authorization header is empty
	if len(authHeader) < 10 {
		return "", ""
	}

	// Split the Authorization header into parts
	headerParts := strings.Split(authHeader, " ")

	// Check if the Authorization header has two parts
	if len(headerParts) < 2 {
		return "", ""
	}

	// Split the second part of the Authorization header into parts
	tokenParts := strings.Split(headerParts[1], ":")

	// Check if the second part of the Authorization header has two parts
	if len(tokenParts) < 2 {
		return "", ""
	}

	// Get the userId and enctoken from the token parts
	userId = tokenParts[0]
	enctoken = tokenParts[1]

	// fmt.Println("GetAuthHeaderValues: ", userId, ":", enctoken)

	// Return the userId and enctoken
	return userId, enctoken

}
