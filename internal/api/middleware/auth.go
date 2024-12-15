package middleware

import (
	"github.com/christo-andrew/haven/pkg/auth"
	"github.com/gin-gonic/gin"
)

func WithAuthUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(401, gin.H{"message": "Unauthorized"})
			return
		}

		claims, err := auth.VerifyToken(token[7:])
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"message": "Invalid token"})
			return
		}

		c.Set("user", claims)
		c.Next()
	}
}
