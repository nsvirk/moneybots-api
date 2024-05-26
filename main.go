// Package implements the Moneybots API.
package main

import (
	"github.com/nsvirk/moneybotsapi/api"
)

// main is the entry point for the application
func main() {

	// Load the environment configuration
	api.LoadEnvConfig()

	// start the server
	api.StartServer()

}
