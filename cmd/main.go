package main

import (
	"log"
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

	username := os.Getenv("LD_USERNAME")
	password := os.Getenv("LD_PASSWORD")

	cookie := services.Login(username, password)

	err := os.Setenv("LD_COOKIE", cookie)

	if err != nil {
		log.Fatal(err)
	}

	router := routes.Init()

	http.ListenAndServe(":8080", router)
}
