package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/christo-andrew/haven/internal/api/requests"
	"github.com/christo-andrew/haven/internal/api/responses"
	"github.com/christo-andrew/haven/internal/models"
	utils "github.com/christo-andrew/haven/pkg/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetAllUsersHandler GetUsers godoc
// @Summary Get all users
// @Description Retrieve all users
// @Produce json
// @Success 200 {array} responses.UserResponse
// @Router /users [get]
// @Tags users
// @Security AuthToken
// @Param Authorization header string true "Authorization"
func GetAllUsersHandler(c *gin.Context, db *gorm.DB) {
	c.JSON(http.StatusOK, getUsers(db))
}

// GetUserHandler GetUser godoc
// @Summary Get a user
// @Description Retrieve a user
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} responses.UserResponse
// @Router /users/{id} [get]
// @Tags users
// @Security AuthToken
// @Param Authorization header string true "Authorization"
func GetUserHandler(c *gin.Context, db *gorm.DB) {
	userId, _ := strconv.Atoi(c.Param("id"))
	user := getUser(userId, db)
	c.JSON(http.StatusOK, user)
}

// CreateUserHandler CreateUser godoc
// @Summary Create a user
// @Description Create a user
// @Accept json
// @Produce json
// @Param user body requests.CreateUserRequest true "Create User Request"
// @Success 201 {object} responses.CreateUserResponse
// @Router /users/create [post]
// @Tags users
func CreateUserHandler(c *gin.Context, db *gorm.DB) {
	var createUserRequest requests.CreateUserRequest
	err := c.ShouldBindJSON(&createUserRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{Message: err.Error()})
		return
	}
	user := createUserRequest.User()
	newUser, err := createUser(user, db)

	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{Message: err.Error()})
		return
	}
	response := responses.CreateUserResponse{}.FromUser(newUser)
	c.JSON(http.StatusCreated, response)
}

func getUser(userID int, db *gorm.DB) *models.User {
	var user models.User
	db.First(&user, userID)
	return &user
}

func getUsers(db *gorm.DB) []models.User {
	var users []models.User
	db.Find(&users)
	return users
}

func createUser(user *models.User, db *gorm.DB) (*models.User, error) {
	// Check if the user already exists
	var existingUser models.User
	db.Where("email = ?", user.Email).First(&existingUser)

	if existingUser != (models.User{}) {
		return &models.User{}, errors.New("user already exists")
	}
	// Hash the password and create the user
	_, hashedPassword := utils.HashPassword(user.Password)
	user.Password = hashedPassword
	db.Create(&user)

	return user, nil
}
