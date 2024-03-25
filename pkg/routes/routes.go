package routes

import (
	"github.com/dronept/go-stremio-legendasdivx/pkg/handlers"
	"github.com/dronept/go-stremio-legendasdivx/pkg/middleware"
	"github.com/dronept/go-stremio-legendasdivx/pkg/services"
	legendasdivx "github.com/dronept/go-stremio-legendasdivx/pkg/services/legendas_divx"
	"github.com/gin-gonic/gin"
)

func CreateRouter(services *services.Services, cache *legendasdivx.SubtitleCache) *gin.Engine {
	router := gin.Default()
	router.LoadHTMLFiles("templates/configure.tmpl")

	h := handlers.NewHandlers(services, cache)

	// Configure CORS
	router.Use(middleware.CORSMiddleware())
	userRoutes := router.Group("/:config")

	router.GET("/manifest.json", h.GetManifestHandler)
	router.GET("/configure", h.GetConfigureHandler)
	router.POST("/configure", h.PostConfigureHandler)

	userRoutes.GET("/configure", h.GetConfigureHandler)
	userRoutes.POST("/configure", h.PostConfigureHandler)
	userRoutes.GET("/manifest.json", h.GetManifestHandler)
	userRoutes.GET("/subtitles/:type/:id/*metadata", h.GetSubtitlesHandler)
	userRoutes.GET("/d/:lid/:id/sub.vtt", h.DownloadSubtitlesHandler)

	return router
}
