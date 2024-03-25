package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/dronept/go-stremio-legendasdivx/configs"
	"github.com/dronept/go-stremio-legendasdivx/pkg/services"
	legendasdivx "github.com/dronept/go-stremio-legendasdivx/pkg/services/legendas_divx"
	"github.com/gin-gonic/gin"
)

type SubtitleResponse struct {
	Id       string `json:"id"`
	Url      string `json:"url"`
	Language string `json:"lang"`
}

func getSubtitlesHandler(services *services.Services, cache *legendasdivx.SubtitleCache) func(c *gin.Context) {
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

		ldSubtitles, err := services.LegendasDivx.GetSubtitles(imdbId, cookie)

		if err != nil && err.Error() == "Login failed" {
			getCookie(c, true, services)
			getSubtitlesHandler(services, cache)(c)
			return
		}

		for i, subtitle := range ldSubtitles {
			name := subtitle.Name
			id := subtitle.Id

			if name == "" {
				name = "subtitle"
			}

			if id == "" {
				id = strconv.Itoa(i)
			}

			// subtitle.DownloadUrl,
			downloadUrl := fmt.Sprintf("%s/%s/d/%s/%d/sub.vtt",
				configs.Values.PublicEndpoint,
				url.QueryEscape(c.Param("config")),
				subtitle.DownloadUrl,
				i,
			)

			cache.Set(subtitle.DownloadUrl, strconv.Itoa(i), subtitle.Name)

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
