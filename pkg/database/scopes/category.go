package scopes

import (
	"github.com/christo-andrew/haven/internal/models"
	"gorm.io/gorm"
)

func GetTransactionTypes(db *gorm.DB) *gorm.DB {
	return db.Where("context = 'accounts' AND context_type = 'transaction_types'")
}

func GetTransactionCategories(db *gorm.DB) *gorm.DB {
	return db.Where("context = 'accounts' AND context_type = 'transaction_categories'")
}

func GetOrCreateTransactionType(name string, db *gorm.DB) *models.Category {
	var transactionType models.Category
	db.Scopes(GetTransactionTypes).Where("name = ?", name).First(&transactionType)
	if transactionType.Name == "" {
		transactionType = models.Category{Name: name, Context: "accounts", ContextType: "transaction_types"}
		db.Create(&transactionType)
	}
	return &transactionType
}

func GetOrCreateTransactionCategory(name string, db *gorm.DB) *models.Category {
	var transactionCategory models.Category
	db.Scopes(GetTransactionCategories).Where("name = ?", name).First(&transactionCategory)
	if transactionCategory.Name == "" {
		transactionCategory = models.Category{Name: name, Context: "accounts", ContextType: "transaction_categories"}
		db.Create(&transactionCategory)
	}
	return &transactionCategory
}

func GetCategoriesByContextAndContextType(context string, contextType string, db *gorm.DB) *gorm.DB {
	return db.Where("context = ? AND context_type = ?", context, contextType)
}
