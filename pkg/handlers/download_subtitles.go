package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/asticode/go-astisub"
	"github.com/dronept/go-stremio-legendasdivx/pkg/services"
	legendasdivx "github.com/dronept/go-stremio-legendasdivx/pkg/services/legendas_divx"
	"github.com/gin-gonic/gin"
	"github.com/saintfish/chardet"
)

func scoredMatch(a, b string) int {
	// This function should split the strings by spaces, and then compare each word and add a +1 to score if they match

	score := 0
	reSplit := regexp.MustCompile(`(\w{2,})`)

	words := reSplit.FindAllString(strings.ToLower(b), -1)

	a = strings.ToLower(a)

	for _, word := range words {
		if strings.Contains(a, word) {
			score++
		}
	}

	return score
}

func downloadSubtitlesHandler(services *services.Services, cache *legendasdivx.SubtitleCache) func(c *gin.Context) {
	return func(c *gin.Context) {
		// get :lid and :name from params
		lid := c.Param("lid")
		id := c.Param("id")

		name, hasCache := cache.Get(lid, id)

		if !hasCache {
			c.String(http.StatusNotFound, "")
			return
		}

		// Download
		files := services.LegendasDivx.Download(lid, getCookie(c, false, services))

		lastScore := -1
		bestScoreIndex := 0

		decodedName, _ := url.QueryUnescape(name.(string))

		fmt.Println("- Matching: ", decodedName)

		// find name inside files, else return first
		for i, fname := range files {
			if score := scoredMatch(path.Base(fname), decodedName); score > lastScore {
				lastScore = score
				bestScoreIndex = i
			}
		}

		fmt.Println("- Best match: ", files[bestScoreIndex])

		// Read file
		sub, err := decode(files[bestScoreIndex])

		// Read file
		// sub, err := astisub.OpenFile(files[bestScoreIndex])

		if err != nil {
			// Handle the error here
			fmt.Println("Error opening subtitle file:", err)
			c.String(http.StatusInternalServerError, "")
			return
		}

		c.Header("Content-type", "text/vtt")

		// Convert to VTT
		var buff = &bytes.Buffer{}
		err = sub.WriteToWebVTT(buff)

		if err != nil {
			// Handle the error here
			fmt.Println("Error converting subtitle file to VTT:", err)
			c.String(http.StatusInternalServerError, "")
			return
		}

		// Write file
		c.String(http.StatusOK, buff.String())
	}
}

func decode(filename string) (o *astisub.Subtitles, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Detect encoding
	detector := chardet.NewTextDetector()
	buf := make([]byte, 1024)
	n, _ := file.Read(buf)
	result, _ := detector.DetectBest(buf[:n])

	fmt.Println("Detected charset: ", result.Charset)

	// if result.Charset == "ISO-8859-1" {
	// 	decodingReader := transform.NewReader(file, charmap.ISO8859_1.NewDecoder())

	// 	return astisub.ReadFromSRT(decodingReader)
	// }

	return astisub.ReadFromSRT(file)
}
