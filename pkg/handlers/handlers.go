package handlers

import (
	"github.com/dronept/go-stremio-legendasdivx/pkg/services"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	GetManifestHandler       func(c *gin.Context)
	GetConfigureHandler      func(c *gin.Context)
	PostConfigureHandler     func(c *gin.Context)
	GetSubtitlesHandler      func(c *gin.Context)
	DownloadSubtitlesHandler func(c *gin.Context)
}

func NewHandlers(s *services.Services) *Handlers {
	return &Handlers{
		GetManifestHandler:       getManifestHandler,
		GetConfigureHandler:      getConfigureHandler,
		PostConfigureHandler:     postConfigureHandler,
		GetSubtitlesHandler:      getSubtitlesHandler(s),
		DownloadSubtitlesHandler: downloadSubtitlesHandler(s),
	}
}
