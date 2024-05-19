package api

import (
	"github.com/labstack/echo/v4"
)

// main is the entry point for the application
func SetupRoutes(e *echo.Echo) {

	// index route
	e.GET("/", IndexHandler)

	// user routes
	e.POST("/user/create", UserCreateHandler)
	e.PATCH("/user/update", UserUpdateHandler)
	e.POST("/user/delete", UserDeleteHandler)
	e.POST("/user/login", UserLoginHandler)
	e.POST("/user/logout", UserLogutHandler)
	e.POST("/user/profile", UserProfileHandler)

}
