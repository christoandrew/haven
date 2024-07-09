package models

import "gorm.io/gorm"

type RealEstateAccount struct {
	gorm.Model
	Account Account `gorm:"polymorphic:BaseAccount;"`
}
