package handlers

import (
	"encoding/base64"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
)

const MySecret string = "abc&1*~#^2^#s0^=)^^7%b34"

func postConfigureHandler(c *gin.Context) {
	// Get application/x-www-form-urlencoded data
	username := c.PostForm("username")
	password := c.PostForm("password")

	// base64 encode
	encodedCredentials := base64.RawStdEncoding.EncodeToString(
		[]byte(username + ":" + password),
	)

	url, err := url.Parse(os.Getenv("PUBLIC_ENDPOINT"))

	if err != nil {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	hostname := url.Hostname()

	// Redirect to another stremio://url
	c.Redirect(
		302,
		"stremio://"+hostname+"/"+encodedCredentials+"/manifest.json",
	)
}
