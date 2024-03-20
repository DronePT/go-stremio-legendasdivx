package main

import (
	"net/http"

	"github.com/dronept/go-stremio-legendasdivx/pkg/routes"
	"github.com/joho/godotenv"
)

type Subtitle struct {
	DownloadUrl string
	Language    string
	Subtitles   []string
}

func main() {
	godotenv.Load()

	router := routes.Init()
	http.ListenAndServe(":8080", router)
}
