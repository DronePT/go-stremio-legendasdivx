package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/dronept/go-stremio-legendasdivx/configs"
	"github.com/dronept/go-stremio-legendasdivx/pkg/services"
	"github.com/gin-gonic/gin"
)

type SubtitleResponse struct {
	Id       string `json:"id"`
	Url      string `json:"url"`
	Language string `json:"lang"`
}

func getSubtitlesHandler(services *services.Services) func(c *gin.Context) {
	return func(c *gin.Context) {

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

		cookie := getCookie(c, false, services)

		var subtitles []SubtitleResponse

		s, err := services.LegendasDivx.GetSubtitles(imdbId, cookie)

		if err != nil && err.Error() == "Login failed" {
			getCookie(c, true, services)
			getSubtitlesHandler(services)(c)
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

			downloadUrl := fmt.Sprintf("%s/%s/download/%s/%s/sub.vtt",
				configs.Values.PublicEndpoint,
				url.QueryEscape(c.Param("config")),
				subtitle.DownloadUrl,
				url.QueryEscape(name),
			)

			url := fmt.Sprintf("%s%s",
				configs.Values.StremioSubtitleEncoder,
				downloadUrl,
			)

			subtitles = append(subtitles, SubtitleResponse{
				Id:       id,
				Url:      url,
				Language: subtitle.Language,
			})
		}

		c.Header("Cache-Control", "max-age=86400,staleRevalidate=stale-while-revalidate, staleError=stale-if-error, public")

		c.JSON(http.StatusOK, gin.H{
			"subtitles": subtitles,
		})

	}
}
