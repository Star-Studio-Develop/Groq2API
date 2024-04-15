package initialize

import (
	"github.com/gin-gonic/gin"
	"groqai2api/middlewares"
	"groqai2api/router"
)

func InitRouter() *gin.Engine {
	Router := gin.Default()

	Router.Use(middlewares.Cors)

	Router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "https://github.com/Star-Studio-Develop/Groq2API")
	})

	Router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	v1Group := Router.Group("/v1/")
	router.InitRouter(v1Group)

	return Router
}
