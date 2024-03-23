package handlers

import (
	"encoding/base64"
	"strings"

	"github.com/dronept/go-stremio-legendasdivx/pkg/services"
	"github.com/gin-gonic/gin"
)

func parseConfig(config string) (string, string) {
	decodedCredentials, _ := base64.RawStdEncoding.DecodeString(config)
	credentials := strings.Split(string(decodedCredentials), ":")

	return credentials[0], credentials[1]
}

func getCookie(c *gin.Context, relogin bool, services *services.Services) string {
	// Get config from :config param, decode it from base64
	config := c.Param("config")
	username, password := parseConfig(config)

	return services.LegendasDivx.Login(username, password, relogin)
}
