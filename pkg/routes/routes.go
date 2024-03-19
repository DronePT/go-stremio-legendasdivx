package routes

import (
	"net/http"

	"github.com/dronept/go-stremio-legendasdivx/pkg/handlers"
	"github.com/gin-gonic/gin"
)

func Init() *gin.Engine {
	router := gin.Default()

	router.GET("/subtitles/:type/:id.json", handlers.GetSubtitlesHandler)

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	return router
}
