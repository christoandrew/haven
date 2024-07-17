package middleware

import (
	"github.com/gin-gonic/gin"
	//utils "github.com/christo-andrew/haven/pkg/utils"
)

func AuthMiddleware() gin.HandlerFunc {
	// Validate JWT check against user and DB
	return func(c *gin.Context) {
		//token := c.GetHeader("Authorization")
		//if token == "" {
		//	c.AbortWithStatusJSON(401, gin.H{"message": "Unauthorized"})
		//	return
		//}
		//
		//if _, err := utils.VerifyToken(token[7:]); err != nil {
		//	c.AbortWithStatusJSON(401, gin.H{"message": "Invalid token"})
		//	return
		//}

		c.Next()
	}
}
