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

		"resources":  []string{"subtitles"},
		"catalogs":   []interface{}{},
		"types":      []string{"movie", "series"},
		"idPrefixes": []string{"tt"},

		// User data requried for login to legendasdivx
		"behaviorHints": map[string]interface{}{
			"configurable":          true,
			"configurationRequired": !configProvided,
		},
		"config": []interface{}{
			map[string]interface{}{
				"key":      "username",
				"type":     "text",
				"title":    "LegendasDivx Username",
				"required": true,
			},
			map[string]interface{}{
				"key":      "password",
				"type":     "password",
				"title":    "LegendasDivx Password",
				"required": true,
			},
		},
	})
}
