package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func VerifyToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
}

func GetUserIdFromContext(c *gin.Context) int {
	token, _ := c.Get("user")
	userId := token.(*jwt.Token).Claims.(jwt.MapClaims)["UserId"]
	return int(userId.(float64))
}
