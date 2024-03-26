package handlers

import (
	"github.com/dronept/go-stremio-legendasdivx/pkg/services"
	"github.com/dronept/go-stremio-legendasdivx/pkg/utils"
	"github.com/gin-gonic/gin"
)

func getCookie(c *gin.Context, relogin bool, services *services.Services) string {
	// Get config from :config param, decode it from base64
	config := c.Param("config")
	username, password := utils.ParseUserData(config)

	return services.LegendasDivx.Login(username, password, relogin)
}
