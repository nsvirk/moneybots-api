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

	// Loglevel
	switch cfg.EchoLogLevel {
	case "DEBUG":
		e.Logger.SetLevel(log.DEBUG)
	case "INFO":
		e.Logger.SetLevel(log.INFO)
	case "WARN":
		e.Logger.SetLevel(log.WARN)
	case "ERROR":
		e.Logger.SetLevel(log.ERROR)
	case "OFF":
		e.Logger.SetLevel(log.OFF)
	}

	// print the server info
	fmt.Println(fmt.Sprintf("%s %s", cfg.APIName, cfg.APIVersion), "server listening on port", cfg.ServerPort)

	// start the server
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", cfg.ServerPort)))
}
