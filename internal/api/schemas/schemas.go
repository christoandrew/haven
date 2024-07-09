package schemas

import "github.com/christo-andrew/haven/internal/models"

type ITransactionSchema interface {
	Transaction(data map[string]interface{}) *models.Transaction
}

func GetTransactionSchemaFromName(bankName string, account *models.Account) ITransactionSchema {
	switch bankName {
	case "Stanbic":
		return NewStanbicTransactionSchema(account)
	default:
		return nil
	}
}
