package middlewares

import (
	"github.com/gin-gonic/gin"
	"groqai2api/global"
	"strings"
)

func Authorization(c *gin.Context) {
	if global.Authorization != "" {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		if global.Authorization != strings.Replace(authHeader, "Bearer ", "", 1) {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
	}
	c.Next()
}
