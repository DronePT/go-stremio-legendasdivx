package handlers

import (
	"fmt"
	"path"
	"strings"

	"github.com/dronept/go-stremio-legendasdivx/pkg/services"
	"github.com/gin-gonic/gin"
)

func match(a, b string) bool {
	// Replace all spaces with dots, lowercase and compare
	a = strings.ReplaceAll(a, " ", ".")
	b = strings.ReplaceAll(b, " ", ".")

	a = strings.ToLower(a)
	b = strings.ToLower(b)

	fmt.Printf("Comparing:\n%s\n%s\n", a, b)

	return strings.Contains(a, b)
}

func DownloadSubtitlesHandler(c *gin.Context) {
	// get :lid and :name from params
	lid := c.Param("lid")
	name := strings.Split(c.Param("name"), ".srt")[0]

	// Download
	files := services.Download(lid, GetCookie(c))

	var matchedFiled string = ""

	// find name inside files, else return first
	for _, fname := range files {
		if match(path.Base(fname), name) {
			matchedFiled = fname
		}
	}

	if matchedFiled == "" {
		fmt.Println("Using file: #1", files[0])
	} else {
		fmt.Println("Using file: #2", matchedFiled)
	}

	// Send file
	c.File(files[0])
}
