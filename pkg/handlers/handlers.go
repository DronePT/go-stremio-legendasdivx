package handlers

import (
	"encoding/base64"
	"strings"

	"github.com/dronept/go-stremio-legendasdivx/pkg/services"
	"github.com/gin-gonic/gin"
)

func GetCookie(c *gin.Context, relogin bool) string {
	// Get config from :config param, decode it from base64
	config := c.Param("config")
	decodedCredentials, _ := base64.RawStdEncoding.DecodeString(config)
	credentials := strings.Split(string(decodedCredentials), ":")

	cookie := services.Login(credentials[0], credentials[1], relogin)

	return cookie
}
