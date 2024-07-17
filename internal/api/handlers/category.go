package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/christo-andrew/haven/internal/api/requests"
	"github.com/christo-andrew/haven/internal/api/serializers"
	"github.com/christo-andrew/haven/internal/models"
	"github.com/christo-andrew/haven/pkg/database/scopes"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetAllCategoriesHandler GetCategories godoc
// @Summary Get all categories
// @Description Retrieve all categories
// @Produce json
// @Success 200 {array} responses.CategoryResponse
// @Router /categories [get]
// @Tags categories
// @Security AuthToken
// @Param Authorization header string true "Authorization"
func GetAllCategoriesHandler(c *gin.Context, db *gorm.DB) {
	var categories []models.Category
	db.Model(&models.Category{}).Find(&categories)
	serializer := serializers.NewCategorySerializer(categories, true)
	response, err := serializer.Serialize()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetCategoryHandler GetCategory godoc
// @Summary Get a category
// @Description Retrieve a category
// @Produce json
// @Success 200 {object} responses.CategoryResponse
// @Router /categories/{id} [get]
// @Param id path int true "Category ID"
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Tags categories
// @Security AuthToken
// @Param Authorization header string true "Authorization"
func GetCategoryHandler(c *gin.Context, db *gorm.DB) {
	var category models.Category
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	category, err = getCategory(id, db)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	serializer := serializers.NewCategorySerializer(category, false)
	response, err := serializer.Serialize()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}

func getCategory(id int, db *gorm.DB) (models.Category, error) {
	var category models.Category
	db.First(&category, id)

	if category.ID == 0 {
		return category, errors.New("category not found")
	}

	return category, nil
}

// CreateCategoryHandler CreateCategory godoc
// @Summary Create a category
// @Description Create a category
// @Accept json
// @Produce json
// @Param category body requests.CreateCategoryRequest true "Create Category Request"
// @Success 201 {object} responses.CategoryResponse
// @Failure 400 {object} responses.ErrorResponse
// @Router /categories/create [post]
// @Tags categories
// @Security AuthToken
// @Param Authorization header string true "Authorization"
func CreateCategoryHandler(c *gin.Context, db *gorm.DB) {
	var createCategoryRequest requests.CreateCategoryRequest
	err := c.ShouldBindJSON(&createCategoryRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	category, err := createCategory(createCategoryRequest, db)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := serializers.NewCategorySerializer(*category, false).Serialize()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, response)
}

func createCategory(createCategoryRequest requests.CreateCategoryRequest, db *gorm.DB) (*models.Category, error) {
	category := createCategoryRequest.Category()
	result := db.Create(category)
	if result.Error != nil {
		return category, result.Error
	}
	return category, nil
}

// GetCategoryByContextAndContextTypeHandler GetTypes godoc
// @Summary Get all category types
// @Description Retrieve all category types
// @Produce json
// @Success 200 {array} responses.CategoryResponse
// @Param context query string false "Context"
// @Param context_type query string false "Context Type"
// @Router /categories/types [get]
// @Tags categories
// @Security AuthToken
// @Param Authorization header string true "Authorization"
func GetCategoryByContextAndContextTypeHandler(c *gin.Context, db *gorm.DB) {
	context := c.Query("context")
	contextType := c.Query("context_type")
	var categories []models.Category
	scopes.GetCategoriesByContextAndContextType(context, contextType, db).Find(&categories)
	serializer := serializers.NewCategorySerializer(categories, true)
	response, err := serializer.Serialize()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}
