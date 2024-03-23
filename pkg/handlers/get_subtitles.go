package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/dronept/go-stremio-legendasdivx/pkg/services"
	"github.com/gin-gonic/gin"
)

type SubtitleResponse struct {
	Id       string `json:"id"`
	Url      string `json:"url"`
	Language string `json:"lang"`
}

func GetSubtitlesHandler(c *gin.Context) {
	// Get imdb id from request parmams
	mediaType := c.Param("type")
	imdbId := c.Param("id")
	// extra := c.Param("extra.json")

	if mediaType != "movie" && mediaType != "series" {
		c.JSON(http.StatusOK, gin.H{
			"subtitles": []any{},
		})
		return
	}

	cookie := GetCookie(c, false)

	var subtitles []SubtitleResponse

	s, err := services.FetchSubtitles(imdbId, cookie)

	if err != nil && err.Error() == "Login failed" {
		GetCookie(c, true)
		GetSubtitlesHandler(c)
		return
	}

	for i, subtitle := range s {
		name := subtitle.Name
		id := subtitle.Id

		if name == "" {
			name = "subtitle"
		}

		if id == "" {
			id = strconv.Itoa(i)
		}

		downloadUrl := fmt.Sprintf("%s/%s/download/%s/%s",
			os.Getenv("PUBLIC_ENDPOINT"),
			c.Param("config"),
			subtitle.DownloadUrl,
			name,
		)

		url := fmt.Sprintf("%s%s",
			os.Getenv("STREMIO_SUBTITLE_PREFIX"),
			url.QueryEscape(downloadUrl),
		)

		subtitles = append(subtitles, SubtitleResponse{
			Id:       id,
			Url:      url,
			Language: subtitle.Language,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"subtitles": subtitles,
	})
}
