package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func getManifestHandler(c *gin.Context) {
	// Get config parameter
	_, configProvided := c.Params.Get("config")

	// Respond with stremio manifest json for subtitles addon
	c.JSON(http.StatusOK, gin.H{
		"id":          "com.legendasdivx",
		"version":     "1.0.0",
		"name":        "LegendasDivx",
		"description": "LegendasDivx subtitles addon for Stremio",

		"resources": []map[string]any{
			{
				"name":  "subtitles",
				"types": []string{"movie", "series"},
			},
		},
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
