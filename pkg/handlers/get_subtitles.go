package handlers

import (
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
	metaType := c.Param("type")
	id := strings.Replace(c.Param("id.json"), ".json", "", -1)
	// extra := c.Param("extra.json")

	fmt.Println("Extra stuff:", id)

	if metaType != "movies" {
		c.JSON(http.StatusOK, gin.H{
			"subtitles": []any{},
		})
		return
	}

	var subtitles []SubtitleResponse

	s := services.FetchSubtitles(id, os.Getenv("LD_COOKIE"))

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
