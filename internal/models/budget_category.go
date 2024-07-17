package models

import "gorm.io/gorm"

type BudgetCategory struct {
	gorm.Model
	BudgetID   int       `json:"budget_id"`
	Budget     *Budget   `json:"budget"`
	CategoryID int       `json:"category_id"`
	Category   *Category `json:"category"`
}
