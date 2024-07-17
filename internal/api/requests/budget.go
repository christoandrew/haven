package requests

import (
	"github.com/christo-andrew/haven/internal/models"
	"time"
)

type CreateBudgetRequest struct {
	Name             string  `json:"name" binding:"required"`
	Description      string  `json:"description"`
	Amount           float64 `json:"amount" binding:"required"`
	UserId           uint    `json:"user_id" binding:"required"`
	BudgetCategoryID uint    `json:"category_id" binding:"required"`
	StartDate        string  `json:"start_date" binding:"required"`
	EndDate          string  `json:"end_date" binding:"required"`
}

func (createBudgetRequest *CreateBudgetRequest) Budget() *models.Budget {
	startDate, _ := time.Parse(createBudgetRequest.StartDate, "2006-01-02")
	endDate, _ := time.Parse(createBudgetRequest.EndDate, "2006-01-02")
	return &models.Budget{
		Name:        createBudgetRequest.Name,
		Description: createBudgetRequest.Description,
		Amount:      createBudgetRequest.Amount,
		UserId:      createBudgetRequest.UserId,
		CategoryID:  createBudgetRequest.BudgetCategoryID,
		StartDate:   startDate,
		EndDate:     endDate,
	}
}
