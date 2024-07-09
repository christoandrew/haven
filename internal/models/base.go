package models

import (
	database "github.com/christo-andrew/haven/pkg/database"
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
	)
	if err != nil {
		panic(err)
	}
}

func DB() *gorm.DB {
	return database.GetDB(database.DatabaseConfig{})
}
