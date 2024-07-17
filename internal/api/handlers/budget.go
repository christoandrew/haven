package handlers

import (
	"github.com/christo-andrew/haven/internal/api/requests"
	"github.com/christo-andrew/haven/internal/api/responses"
	"github.com/christo-andrew/haven/internal/api/serializers"
	"github.com/christo-andrew/haven/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

// CreateBudgetHandler CreateBudget godoc
// @Summary Create a budget
// @Description Create a budget
// @Param budget body requests.CreateBudgetRequest true "Budget"
// @Accept json
// @Produce json
// @Success 201 {object} responses.BudgetResponse
// @Router /budgets/create [post]
// @Tags budgets
// @Security AuthToken
// @Param Authorization header string true "Authorization"
func CreateBudgetHandler(c *gin.Context, db *gorm.DB) {
	var createBudgetRequest requests.CreateBudgetRequest
	if err := c.ShouldBindJSON(&createBudgetRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	budget, err := createBudget(createBudgetRequest, db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, responses.BudgetResponse{}.FromBudget(*budget))
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

// GetAllBudgetsHandler GetBudgets godoc
// @Summary Get all budgets
// @Description Retrieve all budgets
// @Produce json
// @Success 200 {array} responses.BudgetResponse
// @Router /budgets [get]
// @Tags budgets
// @Security AuthToken
// @Param Authorization header string true "Authorization"
func GetAllBudgetsHandler(c *gin.Context, db *gorm.DB) {
	var budgets []models.Budget
	db.Find(&budgets)
	result, err := serializers.NewBudgetSerializer(budgets, true).Serialize()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func createBudget(createBudgetRequest requests.CreateBudgetRequest, db *gorm.DB) (*models.Budget, error) {
	budget := createBudgetRequest.Budget()
	result := db.Create(&budget)
	if result.Error != nil {
		return budget, result.Error
	}

	return budget, nil
}
