package models

import (
	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password" gorm:"type:varchar(256)"`
}

func (user *User) GenerateTokenWithClaims() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"id":    user.ID,
	})

	return token.SignedString([]byte("secret"))
}

func (user *User) GetFullName() string {
	return user.FirstName + " " + user.LastName
}
