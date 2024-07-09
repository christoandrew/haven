package models

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Context     string `json:"context"`
	ContextType string `json:"context_type"`
}
