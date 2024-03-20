package handlers

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strings"

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
	id := c.Param("id")
	// extra := c.Param("extra.json")

	if mediaType != "movie" {
		c.JSON(http.StatusOK, gin.H{
			"subtitles": []any{},
		})
		return
	}

	// Get config from :config param, decode it from base64
	config := c.Param("config")
	decodedCredentials, _ := base64.RawStdEncoding.DecodeString(config)
	credentials := strings.Split(string(decodedCredentials), ":")

	cookie := services.Login(credentials[0], credentials[1])

	var subtitles []SubtitleResponse

	s := services.FetchSubtitles(id, cookie)

	for i, subtitle := range s {
		sid := subtitle.Name
		if sid == "" {
			sid = fmt.Sprint(i)
		}

		url := fmt.Sprintf("%s%s/download/%s/%s.srt",
			os.Getenv("STREMIO_SUBTITLE_PREFIX"),
			os.Getenv("PUBLIC_ENDPOINT"),
			subtitle.DownloadUrl,
			sid,
		)

		subtitles = append(subtitles, SubtitleResponse{
			Id:       sid,
			Url:      url,
			Language: subtitle.Language,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"subtitles": subtitles,
	})
}
