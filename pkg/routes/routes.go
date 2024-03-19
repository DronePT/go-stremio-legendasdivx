package routes

import (
	"net/http"

	"github.com/dronept/go-stremio-legendasdivx/pkg/handlers"
	"github.com/dronept/go-stremio-legendasdivx/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func Init() *gin.Engine {
	router := gin.Default()

	// Configure CORS
	router.Use(middleware.CORSMiddleware())

	router.GET("/manifest.json", handlers.GetManifestHandler)

	router.GET("/subtitles/:type/:id.json", handlers.GetSubtitlesHandler)

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	return router
}
