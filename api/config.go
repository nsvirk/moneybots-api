// Implements the api package
package api

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var cfg Config

type Config struct {
	APIName          string
	APIVersion       string
	ServerPort       string
	EchoLogLevel     string
	JWTSigningKey    string
	SQLDialectorName string
	PostgresDSN      string
	SqliteDSN        string
	SQLDBLogLevel    string
	RedisURL         string
	TelegramBotToken string
	TelegramChatID   string
}

func LoadEnvConfig() {

	// Load .env file to environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// populate the config struct
	cfg = Config{
		APIName:          GetEnvVar("MB_API_API_NAME", "Moneybots API"),
		APIVersion:       GetEnvVar("MB_API_API_VERSION", "v1"),
		ServerPort:       GetEnvVar("MB_API_SERVER_PORT", "3007"),
		EchoLogLevel:     GetEnvVar("MB_API_ECHO_LOG_LEVEL", "DEBUG"),
		JWTSigningKey:    GetEnvVar("MB_API_JWT_SIGNING_KEY", "@Secr3tK3y#@123$%^&*"),
		SQLDialectorName: GetEnvVar("MB_API_SQL_DIALECTOR_NAME", "postgres"),
		PostgresDSN:      GetEnvVar("MB_API_POSTGRES_DSN", "host=localhost user=mb_user password=notset dbname=moneybots port=5432 sslmode=disable TimeZone=Asia/Kolkata"),
		SqliteDSN:        GetEnvVar("MB_API_SQLITE_DSN", "moneybots.db"),
		SQLDBLogLevel:    GetEnvVar("MB_API_SQL_DB_LOG_LEVEL", "Silent"),
		RedisURL:         GetEnvVar("MB_API_REDIS_URL", "redis://localhost:6379/0"),
		TelegramBotToken: GetEnvVar("MB_API_TELEGRAM_BOT_TOKEN", ""),
		TelegramChatID:   GetEnvVar("MB_API_TELEGRAM_CHAT_ID", ""),
	}

}

// GetEnvVar returns the value of an environment variable
func GetEnvVar(key, fallback string) string {
	var value = os.Getenv(key)
	if value == "" {
		if fallback == "" {
			log.Fatalf("%s env var is not set", key)
		} else {
			log.Printf("%s env var is not set, using default value: %s", key, fallback)
			return fallback
		}
	}
	return value
}
