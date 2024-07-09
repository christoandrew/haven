package models

import (
	"gorm.io/gorm"
)

type IAccount interface {
}

type Account struct {
	gorm.Model
	ID              int    `json:"id"`
	AccountName     string `json:"name"`
	AccountType     string `json:"account_type"`
	Currency        string `json:"currency"`
	UserID          uint   `json:"user_id"`
	Balance         float64
	BaseAccountType string        `json:"base_account_type"`
	BaseAccountID   int           `json:"base_account_id"`
	Transactions    []Transaction `gorm:"foreignKey:AccountID"`
}

func (a *Account) BaseAccount() IAccount {
	var account IAccount
	DB().Table(a.BaseAccountType).First(&account, a.BaseAccountID)
	return account
}
