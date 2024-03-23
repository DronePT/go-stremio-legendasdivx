package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func getConfigureHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "configure.tmpl", nil)
}
