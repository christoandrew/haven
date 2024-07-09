package models

import "gorm.io/gorm"

type Tag struct {
	gorm.Model
	ID   int    `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
}
