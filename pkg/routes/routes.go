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

	userRoutes.GET("/manifest.json", handlers.GetManifestHandler)
	userRoutes.GET("/subtitles/:type/:id/*metadata", handlers.GetSubtitlesHandler)

	return router
}
