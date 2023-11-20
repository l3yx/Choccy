package router

import (
	"embed"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetupFrontendRoutes(r *gin.Engine, assets *embed.FS, index *embed.FS) {
	r.GET("/assets/*filepath", func(c *gin.Context) {
		http.FileServer(http.FS(assets)).ServeHTTP(c.Writer, c.Request)
	})
	r.GET("/", func(c *gin.Context) {
		file, _ := index.ReadFile("index.html")
		c.Data(http.StatusOK, "text/html; charset=utf-8", file)
	})
	r.GET("/favicon.ico", func(c *gin.Context) {
		file, _ := index.ReadFile("favicon.ico")
		c.Data(http.StatusOK, "image/x-icon", file)
	})
}
