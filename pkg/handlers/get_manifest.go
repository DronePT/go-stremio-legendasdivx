package handlers

import (
	"net/http"

	"github.com/dronept/go-stremio-legendasdivx/pkg/version"
	"github.com/gin-gonic/gin"
)

func getManifestHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		// Get config parameter
		_, configProvided := c.Params.Get("config")

		// Respond with stremio manifest json for subtitles addon
		c.JSON(http.StatusOK, gin.H{
			"id":          "com.legendasdivx",
			"version":     version.GetVersion(),
			"name":        "LegendasDivx",
			"description": "LegendasDivx subtitles addon for Stremio",

			"resources":  []string{"subtitles"},
			"catalogs":   []interface{}{},
			"types":      []string{"movie", "series"},
			"idPrefixes": []string{"tt"},

			// User data requried for login to legendasdivx
			"behaviorHints": map[string]interface{}{
				"configurable":          true,
				"configurationRequired": !configProvided,
			},
		})
	}
}
