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

	services := services.NewServices()
	router := routes.CreateRouter(services)

	fmt.Printf(
		`
Application version %s
StremioSubtitleEncoder: %s
PublicEndpoint: %s

Server is running on :%s
`,
		version.GetVersion(),
		configs.Values.StremioSubtitleEncoder,
		configs.Values.PublicEndpoint,
		configs.Values.Port,
	)

	err := http.ListenAndServe(":8080", router)

	if err != nil {
		fmt.Println("Error: ", err)
	}
}
