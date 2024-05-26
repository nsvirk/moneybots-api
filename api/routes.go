package api

import (
	"github.com/labstack/echo/v4"
)

// main is the entry point for the application
func SetupRoutes(e *echo.Echo) {

	// index route
	e.GET("/", IndexHandler)

	// user routes
	userRoutes := e.Group("/user")
	userRoutes.POST("/create", UserCreateHandler)
	userRoutes.PATCH("/update", UserUpdateHandler)
	userRoutes.POST("/delete", UserDeleteHandler)
	userRoutes.POST("/login", UserLoginHandler)
	userRoutes.POST("/logout", UserLogutHandler)
	userRoutes.POST("/profile", UserProfileHandler)

	// instrument routes
	instrumentsRoutes := e.Group("/instruments", MWJwtAuthentication())
	instrumentsRoutes.GET("/update", InstrumentsUpdateHandler)
	instrumentsRoutes.GET("/details", InstrumentsDetailsHandler)

	// ticker routes
	tickerRoutes := e.Group("/ticker", MWJwtAuthentication())
	tickerRoutes.POST("/start", TickerStartHandler)
	tickerRoutes.GET("/stop", TickerStopHandler)

}
