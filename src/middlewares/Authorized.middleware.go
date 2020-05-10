package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/enhuizhu/gps-tracking-go-backend/src/models"
)


func Authorized() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !models.IsAuthorized(c) {
			c.JSON(200, gin.H{
				"error": "Unauthorized",
			})
			c.Abort()
			return 
		} else {
			c.Next()
		}		
	}
}