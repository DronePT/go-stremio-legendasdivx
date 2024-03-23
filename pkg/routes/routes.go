package routes

import (
	"github.com/dronept/go-stremio-legendasdivx/pkg/handlers"
	"github.com/dronept/go-stremio-legendasdivx/pkg/middleware"
	"github.com/dronept/go-stremio-legendasdivx/pkg/services"
	"github.com/gin-gonic/gin"
)

func CreateRouter(services *services.Services, appVersion string) *gin.Engine {
	router := gin.Default()
	router.LoadHTMLFiles("templates/configure.tmpl")

	h := handlers.NewHandlers(services, appVersion)

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
	userRoutes.GET("/download/:lid/*name", h.DownloadSubtitlesHandler)

	return router
}
