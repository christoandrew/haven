package handlers

import (
	"github.com/christo-andrew/haven/pkg"
	"net/http"

	"github.com/christo-andrew/haven/internal/api/requests"
	"github.com/christo-andrew/haven/internal/api/responses"
	"github.com/christo-andrew/haven/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// LoginHandler Login godoc
// @Summary Login a user
// @Description Login a user
// @Accept json
// @Produce json
// @Param user body requests.LoginRequest true "Login Request"
// @Success 200 {object} responses.LoginResponse
// @Router /auth/login [post]
// @Tags auth
func LoginHandler(c *gin.Context, db *gorm.DB) {
	var loginRequest requests.LoginRequest
	err := c.BindJSON(&loginRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{Message: err.Error()})
		return
	}
	var user models.User
	user, err = getUserByCredentials(loginRequest.Username, loginRequest.Password, db)
	if err != nil {
		c.JSON(http.StatusUnauthorized, responses.ErrorResponse{Message: err.Error()})
		return
	}
	token, err := user.GenerateToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Message: err.Error()})
		return
	}

	response := responses.LoginResponse{
		Token: token,
	}

	c.JSON(http.StatusOK, response)
}

// The most basic authentication in the world
func getUserByCredentials(username string, password string, db *gorm.DB) (models.User, error) {
	var user models.User
	db.Where("username = ?", username).First(&user)

	if user.ID == 0 {
		return user, pkg.UserNotFoundError()
	}

	if !pkg.ComparePassword(user.Password, password) {
		return user, pkg.InvalidCredentialsError()
	}

	return user, nil
}
