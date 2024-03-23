package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/dronept/go-stremio-legendasdivx/pkg/routes"
	"github.com/dronept/go-stremio-legendasdivx/pkg/services"
	"github.com/joho/godotenv"
)

type Subtitle struct {
	DownloadUrl string
	Language    string
	Subtitles   []string
}

func main() {
	godotenv.Load()

	// STREMIO_SUBTITLE_PREFIX=http://127.0.0.1:11470/subtitles.vtt?from=
	// PUBLIC_ENDPOINT=http://localhost:8080

	fmt.Printf(
		"STREMIO_SUBTITLE_PREFIX: %s\nPUBLIC_ENDPOINT: %s\n\n",
		os.Getenv("STREMIO_SUBTITLE_PREFIX"),
		os.Getenv("PUBLIC_ENDPOINT"),
	)

	services := services.NewServices()
	router := routes.CreateRouter(services)

	fmt.Println("Server running on port :8080")

	err := http.ListenAndServe(":8080", router)

	if err != nil {
		fmt.Println("Error: ", err)
	}
}
