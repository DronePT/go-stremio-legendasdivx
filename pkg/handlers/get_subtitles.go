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
	id := c.Param("id")
	// extra := c.Param("extra.json")

	if mediaType != "movie" {
		c.JSON(http.StatusOK, gin.H{
			"subtitles": []any{},
		})
		return
	}

	cookie := GetCookie(c)

	var subtitles []SubtitleResponse

	s := services.FetchSubtitles(id, cookie)

	for i, subtitle := range s {
		sid := subtitle.Name
		if sid == "" {
			sid = fmt.Sprint(i)
		}

		url := fmt.Sprintf("%s%s/%s/download/%s/%s.srt",
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
