package requests

import (
	"github.com/christo-andrew/haven/internal/models"
	"time"
)

type CreateOrUpdateBudgetRequest struct {
	Name             string  `json:"name" binding:"required"`
	Description      string  `json:"description"`
	Amount           float64 `json:"amount" binding:"required"`
	BudgetCategoryID uint    `json:"category_id" binding:"required"`
	StartDate        string  `json:"start_date" binding:"required"`
	EndDate          string  `json:"end_date" binding:"required"`
}

func (createOrUpdateBudgetRequest *CreateOrUpdateBudgetRequest) Budget() *models.Budget {
	startDate, _ := time.Parse("2006-01-02", createOrUpdateBudgetRequest.StartDate)
	endDate, _ := time.Parse("2006-01-02", createOrUpdateBudgetRequest.EndDate)
	return &models.Budget{
		Name:        createOrUpdateBudgetRequest.Name,
		Description: createOrUpdateBudgetRequest.Description,
		Amount:      createOrUpdateBudgetRequest.Amount,
		CategoryID:  createOrUpdateBudgetRequest.BudgetCategoryID,
		StartDate:   startDate,
		EndDate:     endDate,
	}
}
