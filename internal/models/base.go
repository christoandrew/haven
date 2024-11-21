package models

import (
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&Account{},
		&BankAccount{},
		&Transaction{},
		&CreditCardAccount{},
		&RealEstateAccount{},
		&Category{},
		&User{},
		&Tag{},
		&BudgetCategory{},
		&Budget{},
	)
	if err != nil {
		panic(err)
	}
}
