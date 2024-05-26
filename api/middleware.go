package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// SetupMiddleware sets up the middleware for the application
func SetupMiddleware(e *echo.Echo) {
	e.Use(MWRecover())
	e.Use(MWLogger())
}

// MWLogger middleware logs the request method, uri, and status
func MWLogger() echo.MiddlewareFunc {
	return middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	})
}

// MWRecover middleware recovers from panics anywhere in the chain, logs the panic (and a stack trace), and returns a HTTP 500 (Internal Server Error) status if possible.
// It is recommended to use this middleware at the top of the chain.
func MWRecover() echo.MiddlewareFunc {
	return middleware.Recover()
}

// MWKeyAuthentication middleware authenticates the user
func MWKeyAuthentication() echo.MiddlewareFunc {
	return middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup:  "header:" + echo.HeaderAuthorization,
		AuthScheme: "Bearer",

		Validator: func(key string, c echo.Context) (bool, error) {
			// fmt.Println("authMiddleware key:", key)
			keyParts := strings.Split(key, ":")
			userId := keyParts[0]
			enctoken := keyParts[1]
			result, err := isUserAuthorized(userId, enctoken)
			return result, err
		},

		ErrorHandler: func(err error, c echo.Context) error {
			statusCode := http.StatusUnauthorized
			message := fmt.Sprintf("Unauthorized user: %v", err.Error())
			return SendError(c, statusCode, message)
		},
	})

}

// isUserAuthorized compares userId and enctoken with the db values
func isUserAuthorized(userId, enctoken string) (bool, error) {

	// connect to the database
	db := ConnectToDB()

	// check if user exists
	var dbUser = UserModel{}
	result := db.First(&dbUser, &UserModel{
		UserId:   userId,
		Enctoken: enctoken,
	})

	userAuthorized := result.RowsAffected == 1

	return userAuthorized, nil
}
