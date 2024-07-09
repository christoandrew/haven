package requests

import "github.com/christo-andrew/haven/internal/models"

type CreateCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Context     string `json:"context"`
	ContextType string `json:"context_type"`
}

func (c CreateCategoryRequest) Category() *models.Category {
	return &models.Category{
		Name:        c.Name,
		Description: c.Description,
		Context:     c.Context,
		ContextType: c.ContextType,
	}
}
