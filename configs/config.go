package configs

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	StremioSubtitleEncoder string
	PublicEndpoint         string
}

var Values *Config

func InitConfig() {
	godotenv.Load()

	Values = &Config{
		StremioSubtitleEncoder: os.Getenv("STREMIO_SUBTITLE_PREFIX"),
		PublicEndpoint:         os.Getenv("PUBLIC_ENDPOINT"),
	}
}
