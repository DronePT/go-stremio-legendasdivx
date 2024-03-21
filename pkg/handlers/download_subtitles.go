package handlers

import (
	"fmt"
	"net/url"
	"path"
	"regexp"
	"strings"

	"github.com/dronept/go-stremio-legendasdivx/pkg/services"
	"github.com/gin-gonic/gin"
)

func scoredMatch(a, b string) int {
	// This function should split the strings by spaces, and then compare each word and add a +1 to score if they match

	score := 0
	reSplit := regexp.MustCompile(`(?mi)\W+`)

	words := reSplit.Split(b, -1)

	for _, word := range words {
		if strings.Contains(a, word) {
			score++
		}
	}

	return score
}

func DownloadSubtitlesHandler(c *gin.Context) {
	// get :lid and :name from params
	lid := c.Param("lid")
	name := strings.Split(c.Param("name"), ".srt")[0]

	// Download
	files := services.Download(lid, GetCookie(c, false))

	lastScore := -1
	bestScoreIndex := 0

	decodedName, _ := url.QueryUnescape(name)

	fmt.Println("- Matching: ", decodedName)

	// find name inside files, else return first
	for i, fname := range files {
		if score := scoredMatch(path.Base(fname), decodedName); score > lastScore {
			lastScore = score
			bestScoreIndex = i
		}
	}

	fmt.Println("- Best match: ", files[bestScoreIndex])

	// Send file
	c.File(files[bestScoreIndex])
}
