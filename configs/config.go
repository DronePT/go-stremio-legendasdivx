package configs

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                   string
	StremioSubtitleEncoder string
	PublicEndpoint         string
}

var Values *Config

func InitConfig() {
	godotenv.Load()

	Values = &Config{
		Port:                   getEnv("PORT", "8080"),
		StremioSubtitleEncoder: getEnv("STREMIO_SUBTITLE_PREFIX", ""),
		PublicEndpoint:         getEnv("PUBLIC_ENDPOINT", "http://localhost:8080"),
	}
}

func getEnv(env, defaultValue string) string {
	value := os.Getenv(env)

	if value == "" {
		return defaultValue
	}

	return value
}
