package handlers

import (
	"fmt"
	"net/http"
	"os"

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

	if mediaType != "movie" {
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
		sid := subtitle.Name
		if sid == "" {
			sid = fmt.Sprint(i)
		}

		url := fmt.Sprintf("%s%s/%s/download/%s/%s.vtt",
			os.Getenv("STREMIO_SUBTITLE_PREFIX"),
			os.Getenv("PUBLIC_ENDPOINT"),
			c.Param("config"),
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
