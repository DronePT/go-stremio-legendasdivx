package main

import (
	"fmt"
	"net/http"

	"github.com/dronept/go-stremio-legendasdivx/configs"
	"github.com/dronept/go-stremio-legendasdivx/pkg/routes"
	"github.com/dronept/go-stremio-legendasdivx/pkg/services"
	"github.com/dronept/go-stremio-legendasdivx/pkg/version"
)

type Subtitle struct {
	DownloadUrl string
	Language    string
	Subtitles   []string
}

func main() {
	configs.InitConfig()

	fmt.Printf(
		"STREMIO_SUBTITLE_PREFIX: %s\nPUBLIC_ENDPOINT: %s\n\n",
		configs.Values.StremioSubtitleEncoder,
		configs.Values.PublicEndpoint,
	)

	services := services.NewServices()
	router := routes.CreateRouter(services)

	fmt.Printf(
		"Application version %s\nServer is running on :%s",
		version.GetVersion(),
		configs.Values.Port,
	)

	err := http.ListenAndServe(":8080", router)

	if err != nil {
		fmt.Println("Error: ", err)
	}
}
