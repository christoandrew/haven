package models

import (
	"gorm.io/gorm"
	"time"
)

type Budget struct {
	gorm.Model
	Id          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	UserId      uint      `json:"user_id"`
	User        User      `json:"user"`
	Category    Category  `json:"category"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	CategoryID  uint      `json:"category_id"`
}
