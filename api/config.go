// Implements the api package
package api

import (
	"log"
	"os"
)

// API_NAME is the name of the application
const API_NAME = "Moneybots API"

// API_VERSION is the version of the API
const API_VERSION = "v1"

// GetServerPort returns the port the server listens on
func GetServerPort() string {
	var port = os.Getenv("MB_API_SERVER_PORT")
	if port == "" {
		log.Fatalln("MB_API_SERVER_PORT env is not set")
	}
	return port
}

// GetPostgresDsn returns the postgres dsn
func GetPostgresDsn() string {
	var dsn = os.Getenv("MB_API_POSTGRES_DSN")
	if dsn == "" {
		log.Fatalln("MB_API_POSTGRES_DSN env is not set")
	}
	return dsn
}

// GetSqliteDsn returns the sqlite dsn
func GetSqliteDsn() string {
	var dsn = os.Getenv("MB_API_SQLITE_DSN")
	if dsn == "" {
		log.Fatalln("MB_API_SQLITE_DSN env is not set")
	}
	return dsn
}
