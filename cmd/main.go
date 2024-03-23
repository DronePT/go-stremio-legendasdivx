package main

import (
	"fmt"
	"net/http"

	"github.com/dronept/go-stremio-legendasdivx/configs"
	"github.com/dronept/go-stremio-legendasdivx/pkg/routes"
	"github.com/dronept/go-stremio-legendasdivx/pkg/services"
)

type Subtitle struct {
	DownloadUrl string
	Language    string
	Subtitles   []string
}

const AppVersion = "1.0.1"

func main() {
	configs.InitConfig()

	fmt.Printf(
		"STREMIO_SUBTITLE_PREFIX: %s\nPUBLIC_ENDPOINT: %s\n\n",
		configs.Values.StremioSubtitleEncoder,
		configs.Values.PublicEndpoint,
	)

	services := services.NewServices()
	router := routes.CreateRouter(services, AppVersion)

	fmt.Println("Server running on port :8080")

	err := http.ListenAndServe(":8080", router)

	if err != nil {
		fmt.Println("Error: ", err)
	}
}
