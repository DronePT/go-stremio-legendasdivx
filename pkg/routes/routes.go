package routes

import (
	"github.com/dronept/go-stremio-legendasdivx/pkg/handlers"
	"github.com/dronept/go-stremio-legendasdivx/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func Init() *gin.Engine {
	router := gin.Default()
	router.LoadHTMLFiles("templates/configure.tmpl")

	// Configure CORS
	router.Use(middleware.CORSMiddleware())

	router.GET("/manifest.json", handlers.GetManifestHandler)
	router.GET("/configure", handlers.GetConfigureHandler)
	router.POST("/configure", handlers.PostConfigureHandler)

	userRoutes := router.Group("/:config")

	userRoutes.GET("/configure", handlers.GetConfigureHandler)
	userRoutes.POST("/configure", handlers.PostConfigureHandler)
	userRoutes.GET("/manifest.json", handlers.GetManifestHandler)
	userRoutes.GET("/subtitles/:type/:id/*metadata", handlers.GetSubtitlesHandler)
	userRoutes.GET("/download/:lid/*name", handlers.DownloadSubtitlesHandler)

	return router
}
