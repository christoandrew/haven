package models

import "gorm.io/gorm"

type BankAccount struct {
	gorm.Model
	Account Account `gorm:"polymorphic:BaseAccount;"`
}
