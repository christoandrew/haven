package handlers

import (
	"github.com/christo-andrew/haven/internal/api/requests"
	"github.com/christo-andrew/haven/internal/api/responses"
	"github.com/christo-andrew/haven/internal/api/serializers"
	"github.com/christo-andrew/haven/internal/models"
	"github.com/christo-andrew/haven/pkg/auth"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

// CreateBudgetHandler CreateBudget godoc
// @Summary Create a budget
// @Description Create a budget
// @Param budget body requests.CreateOrUpdateBudgetRequest true "Budget"
// @Accept json
// @Produce json
// @Success 201 {object} responses.BudgetResponse
// @Router /budgets/create [post]
// @Tags budgets
// @Security AuthToken
// @Param Authorization header string true "Authorization"
func CreateBudgetHandler(c *gin.Context, db *gorm.DB) {
	userId := auth.GetUserIdFromContext(c)
	var createOrUpdateBudgetRequest requests.CreateOrUpdateBudgetRequest
	if err := c.ShouldBindJSON(&createOrUpdateBudgetRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	budget := createOrUpdateBudgetRequest.Budget()
	budget.UserId = uint(userId)
	budget, err := createBudget(budget, db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, responses.BudgetResponse{}.FromBudget(*budget))
}

// UpdateBudgetHandler UpdateBudget godoc
// @Summary Update a budget
// @Description Update a budget
// @Param id path int true "Budget ID"
// @Param budget body requests.CreateOrUpdateBudgetRequest true "Budget"
// @Accept json
// @Produce json
// @Success 200 {object} responses.BudgetResponse
// @Router /budgets/{id}/update [put]
// @Tags budgets
// @Security AuthToken
// @Param Authorization header string true "Authorization"
func UpdateBudgetHandler(c *gin.Context, db *gorm.DB) {
	var updateBudgetRequest requests.CreateOrUpdateBudgetRequest
	if err := c.ShouldBindJSON(&updateBudgetRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var budget models.Budget
	id := c.Param("id")
	db.First(&budget, id)
	db.Model(&budget).Updates(updateBudgetRequest.Budget())

	c.JSON(http.StatusOK, responses.BudgetResponse{}.FromBudget(budget))
}

// GetBudgetHandler GetBudget godoc
// @Summary Get a budget
// @Description Retrieve a budget
// @Produce json
// @Success 200 {object} responses.BudgetResponse
// @Router /budgets/{id} [get]
// @Param id path int true "Budget ID"
// @Tags budgets
// @Security AuthToken
// @Param Authorization header string true "Authorization"
func GetBudgetHandler(c *gin.Context, db *gorm.DB) {
	var budget models.Budget
	id := c.Param("id")
	db.First(&budget, id)
	response, err := serializers.NewBudgetSerializer(budget, false).Serialize()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}

// GetBudgetsHandler GetBudgets godoc
// @Summary Get all budgets for a user
// @Description Retrieve all budgets for a user
// @Produce json
// @Success 200 {array} responses.BudgetResponse
// @Router /budgets [get]
// @Tags budgets
// @Security AuthToken
// @Param Authorization header string true "Authorization"
func GetBudgetsHandler(c *gin.Context, db *gorm.DB) {
	userId := auth.GetUserIdFromContext(c)
	var budgets []models.Budget
	db.Where("user_id = ?", userId).Find(&budgets)
	result, err := serializers.NewBudgetSerializer(budgets, true).Serialize()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// AddBudgetTagHandler AddBudgetTag godoc
// @Summary Add a tag to a budget
// @Description Add a tag to a budget
// @Param id path int true "Budget ID"
// @Param tag_id path int true "Tag ID"
// @Produce json
// @Success 200 {object} responses.TagResponse
// @Router /budgets/{id}/tags/ [post]
// @Tags budgets
// @Security AuthToken
// @Param Authorization header string true "Authorization"
func AddBudgetTagHandler(ctx *gin.Context, db *gorm.DB) {
	var budget models.Budget
	id := ctx.Param("id")
	db.First(&budget, id)
	var tag models.Tag
	tagId := ctx.Param("tag_id")
	db.First(&tag, tagId)
	err := db.Model(&budget).Association("Tags").Append(&tag)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, responses.TagResponse{}.FromTag(tag))
}

func RemoveBudgetTagHandler(ctx *gin.Context, db *gorm.DB) error {
	var budget models.Budget
	id := ctx.Param("id")
	db.First(&budget, id)
	var tag models.Tag
	tagId := ctx.Param("tag_id")
	db.First(&tag, tagId)
	err := db.Model(&budget).Association("Tags").Delete(&tag)
	if err != nil {
		return err
	}
	return nil
}

func GetBudgetTagsHandler(ctx *gin.Context, db *gorm.DB) ([]models.Tag, error) {
	var budget models.Budget
	id := ctx.Param("id")
	db.First(&budget, id)
	var tags []models.Tag
	err := db.Model(&budget).Association("Tags").Find(&tags)
	if err != nil {
		return tags, err
	}
	return tags, nil
}

func createBudget(budget *models.Budget, db *gorm.DB) (*models.Budget, error) {
	result := db.Create(&budget)
	if result.Error != nil {
		return budget, result.Error
	}

	return budget, nil
}
