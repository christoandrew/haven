package schemas

import (
	"github.com/christo-andrew/haven/internal/models"
	"gorm.io/gorm"
)

type ITransactionSchema interface {
	Transaction(data map[string]interface{}) *models.Transaction
}

func GetTransactionSchemaFromName(bankName string, account *models.Account, db *gorm.DB) ITransactionSchema {
	switch bankName {
	case "Stanbic":
		return NewStanbicTransactionSchema(account, db)
	default:
		return nil
	}
}
