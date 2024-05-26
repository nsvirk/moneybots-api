package api

import (
	"log"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ConnectToDB connects to the database
func ConnectToDB() *gorm.DB {

	// logLevel
	var logLevel logger.LogLevel

	switch cfg.SQLDBLogLevel {
	case "Silent":
		logLevel = logger.Silent
	case "Error":
		logLevel = logger.Error
	case "Warn":
		logLevel = logger.Warn
	case "Info":
		logLevel = logger.Info
	default:
		logLevel = logger.Silent
	}

	// dialector
	switch cfg.SQLDialectorName {
	case "postgres":
		// connect to the postgres database
		return connectToPostgresDB(logLevel)

	case "sqlite":
		// connect to the sqlite database
		return connectToSqliteDB(logLevel)

	default:
		panic("⇨ invalid SQL_DIALECTOR_NAME")
	}

}

// connectToPostgresDB connects to the postgres database
func connectToPostgresDB(logLevel logger.LogLevel) *gorm.DB {

	// connect to the database
	db, err := gorm.Open(postgres.Open(cfg.PostgresDSN), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})

	if err != nil {
		log.Fatalf("⇨ failed to connect postgres database. %s", err.Error())
	}

	// fmt.Println("⇨ connected to postgres database")
	return db
}

// connectToSqliteDB connects to the sqlite database
func connectToSqliteDB(logLevel logger.LogLevel) *gorm.DB {

	// connect to the database
	db, err := gorm.Open(sqlite.Open(cfg.SqliteDSN), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})

	// check if the connection was successful
	if err != nil {
		log.Fatalf("⇨ failed to connect sqlite database. %s", err.Error())

	}

	// fmt.Println("⇨ connected to sqlite database")
	return db
}

func CloseDB(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("⇨ failed to close the sql database connection. %s", err.Error())
	}
	sqlDB.Close()
}

// MigrateDBSchema migrates the database schema
func MigrateDBSchema() {

	// Connect to the database
	db := ConnectToDB()

	// Migrate the schema
	db.AutoMigrate(&UserModel{})        // UserModel is defined in user_model.go
	db.AutoMigrate(&InstrumentsModel{}) // InstrumentsModel is defined in instruments_model.go

	// Close the database connection
	defer CloseDB(db)
}

// ConnectToRedis connects to the redis database
func ConnectToRedis() *redis.Client {

	// parse the redis url
	opt, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		log.Fatalf("⇨ failed to connect to redis database. %s", err.Error())
	}

	// connect to redis
	return redis.NewClient(opt)

}
