package main

import (
	"fmt"
	"net/http"

	"github.com/dronept/go-stremio-legendasdivx/configs"
	"github.com/dronept/go-stremio-legendasdivx/pkg/routes"
	"github.com/dronept/go-stremio-legendasdivx/pkg/services"
	legendasdivx "github.com/dronept/go-stremio-legendasdivx/pkg/services/legendas_divx"
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
	subtitlesCache := legendasdivx.NewSubtitleCache()
	router := routes.CreateRouter(services, subtitlesCache)

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
