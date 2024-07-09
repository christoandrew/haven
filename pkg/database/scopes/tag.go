package scopes

import (
	"github.com/christo-andrew/haven/internal/models"
	"gorm.io/gorm"
)

func GetOrCreateTransactionTag(name string, db *gorm.DB) *models.Tag {
	var tag models.Tag
	db.Where("name = ?", name).First(&tag)
	if tag.Name == "" {
		tag = models.Tag{Name: name}
		db.Create(&tag)
	}
	return &tag
}
