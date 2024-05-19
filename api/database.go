package api

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ConnectToDB connects to the database
func ConnectToDB() *gorm.DB {

	// ---------------------------------------------------------
	// select the database configuration
	// ---------------------------------------------------------
	// dialector
	// -----------------
	var dialector gorm.Dialector = postgres.Dialector{}
	// var dialector gorm.Dialector = sqlite.Dialector{}

	// logLevel
	// -----------------
	var logLevel logger.LogLevel = logger.Info
	// ---------------------------------------------------------

	// connect to the database
	if dialector.Name() == "postgres" {
		var dsn string = GetPostgresDsn() // GetPostgresDsn is defined in config.go
		return connectToPostgresDB(dsn, logLevel)
	}

	if dialector.Name() == "sqlite" {
		var dsn string = GetSqliteDsn() // GetSqliteDsn is defined in config.go
		return connectToSqliteDB(dsn, logLevel)

	}

	// default to sqlite
	var dsn string = GetSqliteDsn() // GetSqliteDsn is defined in config.go
	logLevel = logger.Silent
	return connectToSqliteDB(dsn, logLevel)
}

// connectToPostgresDB connects to the postgres database
func connectToPostgresDB(dsn string, logLevel logger.LogLevel) *gorm.DB {

	// connect to the database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		log.Fatalln("⇨ Failed to connect to postgres database")
	}

	// defer CloseDB(db)

	fmt.Println("⇨ connected to postgres database")

	// return the database connection
	return db
}

// connectToSqliteDB connects to the sqlite database
func connectToSqliteDB(dsn string, logLevel logger.LogLevel) *gorm.DB {

	// check if the dsn is set
	if dsn == "" {
		panic("⇨ sqlite dsn is not set\n")
	}

	// connect to the database
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})

	// defer CloseDB(db)

	// check if the connection was successful
	if err != nil {
		log.Fatalln("⇨ failed to connect sqlite database")
	}

	fmt.Println("⇨ connected to sqlite database")

	// return the database connection
	return db
}

func CloseDB(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalln("⇨ failed to close the database connection")
	}
	sqlDB.Close()
}

// MigrateDBSchema migrates the database schema
func MigrateDBSchema() {

	// Connect to the database
	db := ConnectToDB()

	// Migrate the schema
	db.AutoMigrate(&UserModel{}) // UserModel is defined in user_model.go

	// Close the database connection

}
