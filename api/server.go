package api

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

// StartServer starts the server
func StartServer() {

	// Migrate the schema
	MigrateDBSchema()

	// instantiate a new echo instance
	e := echo.New()

	// setup middleware in moddleware.go
	SetupMiddleware(e)

	// setup routes in routes.go
	SetupRoutes(e)

	// set echo configurations
	e.HideBanner = true
	e.Logger.SetLevel(log.DEBUG)

	// get the server port
	serverPort := GetServerPort()

	// print the server info
	fmt.Println(fmt.Sprintf("%s %s", API_NAME, API_VERSION), "server listening on port", serverPort)

	// start the server
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", serverPort)))
}
