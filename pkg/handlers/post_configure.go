package handlers

import (
	"encoding/base64"
	"os"

	"github.com/gin-gonic/gin"
)

const MySecret string = "abc&1*~#^2^#s0^=)^^7%b34"

func PostConfigureHandler(c *gin.Context) {
	// Get application/x-www-form-urlencoded data
	username := c.PostForm("username")
	password := c.PostForm("password")

	// base64 encode
	encodedCredentials := base64.RawStdEncoding.EncodeToString(
		[]byte(username + ":" + password),
	)

	// Redirect to another stremio://url
	c.Redirect(
		302,
		"stremio://"+os.Getenv("PUBLIC_ENDPOINT")+"/"+encodedCredentials+"/manifest.json",
	)
}
