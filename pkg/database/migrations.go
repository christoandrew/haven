package database

import (
	"github.com/christo-andrew/haven/internal/models"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&models.Account{},
		&models.BankAccount{},
		&models.Transaction{},
		&models.CreditCardAccount{},
		&models.RealEstateAccount{},
		&models.Category{},
		&models.User{},
		&models.Tag{},
		&models.BudgetCategory{},
		&models.Budget{},
	)
	if err != nil {
		panic(err)
	}
}
