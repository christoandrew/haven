package scopes

import (
	"gorm.io/gorm"
)

func GetRecentTransactions(db *gorm.DB, accountId int, limit int) *gorm.DB {
	return GetTransactionsByAccountId(accountId, db).Limit(limit).Order("date DESC")
}

func GetAccountTransactionsByCategory(accountId int, categoryId int, db *gorm.DB) *gorm.DB {
	return db.Where("account_id = ? AND category_id = ?", accountId, categoryId)
}
