package models

import "gorm.io/gorm"

type CreditCardAccount struct {
	gorm.Model
	Account Account `gorm:"polymorphic:BaseAccount;"`
}
