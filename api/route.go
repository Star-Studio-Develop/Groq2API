package api

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/Star-Studio-Develop/Groq2API/initialize/auth"
	"github.com/Star-Studio-Develop/Groq2API/initialize/stream"
	"github.com/Star-Studio-Develop/Groq2API/initialize/user"
	"github.com/Star-Studio-Develop/Groq2API/initialize/utils"
)

func RegisterRoutes(router *gin.Engine) {
	router.POST("/fetch-jwt", func(c *gin.Context) {
		refreshToken := c.PostForm("refreshToken")
		jwt, err := auth.FetchJWT(refreshToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"jwt": jwt})
	})

	router.GET("/user/profile", func(c *gin.Context) {
		jwt := c.GetHeader("Authorization")
		orgID, err := user.FetchUserProfile(jwt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"org_id": orgID})
	})

	router.POST("/stream/fetch", func(c *gin.Context) {
		var messages []stream.model.Message
		if err := c.BindJSON(&messages); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}
		jwt := c.GetHeader("Authorization")
		orgID := c.PostForm("orgID")
		modelType := c.PostForm("modelType")
		maxTokens := c.GetInt64("maxTokens")

		response, err := stream.FetchStream(jwt, orgID, messages, modelType, maxTokens)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.DataFromReader(http.StatusOK, response.ContentLength, response.Header.Get("Content-Type"), response.Body, nil)
	})

	router.Use(utils.SetCorsHeaders)
}

func main() {
	router := initialize.RegisterRouter()
	RegisterRoutes(router)
	router.Run(":8080")
}
