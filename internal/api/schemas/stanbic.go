package schemas

import (
	"strconv"
	"time"

	"github.com/christo-andrew/haven/internal/models"
	"github.com/christo-andrew/haven/pkg/database/scopes"
	"gorm.io/gorm"
)

type StanbicTransactionSchema struct {
	db      *gorm.DB
	Account *models.Account
}

func NewStanbicTransactionSchema(account *models.Account, db *gorm.DB) *StanbicTransactionSchema {
	return &StanbicTransactionSchema{
		Account: account,
		db:      db,
	}
}

func (schema *StanbicTransactionSchema) Transaction(data map[string]interface{}) *models.Transaction {
	date := data["Date"]
	if date == nil {
		date = time.Now()
	}
	date, err := time.Parse("02/01/2006", date.(string))
	if err != nil {
		date = time.Now()
	}
	description := data["Description"].(string)
	creditAmount, err := strconv.ParseFloat(data["Credit"].(string), 64)
	if err != nil {
		creditAmount = 0
	}
	debitAmount, err := strconv.ParseFloat(data["Debit"].(string), 64)
	if err != nil {
		debitAmount = 0
	}

	amount := creditAmount + debitAmount

	var transactionTypeName string
	if creditAmount > 0 {
		transactionTypeName = "Credit"
	} else if debitAmount > 0 {
		transactionTypeName = "Debit"
	}

	transactionType := scopes.GetOrCreateTransactionType(transactionTypeName, schema.db)
	category := scopes.GetOrCreateTransactionCategory("General", schema.db)

	return &models.Transaction{
		Amount:          amount,
		Date:            date.(time.Time),
		Description:     description,
		TransactionType: *transactionType,
		Currency:        schema.Account.Currency,
		Account:         *schema.Account,
		Category:        *category,
	}
}
