package requests

import (
	"fmt"
	"time"

	"github.com/christo-andrew/haven/internal/models"
	"github.com/christo-andrew/haven/pkg/database/scopes"
	"gorm.io/gorm"
)

type CreateTransactionRequest struct {
	AccountID         int     `json:"account_id"`
	Amount            float64 `json:"amount"`
	Currency          string  `json:"currency"`
	Date              string  `json:"date"`
	Description       string  `json:"description"`
	CategoryID        string  `json:"category_id"`
	TransactionTypeID string  `json:"transaction_type_id"`
	TransactionType   string  `json:"transaction_type"`
	Category          string  `json:"category"`
	DateFormat        string  `json:"date_format"`
}

func (c *CreateTransactionRequest) Transaction(db *gorm.DB) *models.Transaction {
	category := c.GetCategory(db)
	transactionType := c.GetTransactionType(db)

	return &models.Transaction{
		AccountID:         c.AccountID,
		Amount:            c.Amount,
		Currency:          c.Currency,
		Date:              c.FormatDate(),
		Description:       c.Description,
		CategoryID:        category.ID,
		TransactionTypeID: transactionType.ID,
	}

}

func (c *CreateTransactionRequest) GetCategory(db *gorm.DB) *models.Category {
	if c.Category == "" {
		return scopes.GetOrCreateTransactionCategory("General", db)
	}
	return scopes.GetOrCreateTransactionCategory(c.Category, db)
}

func (c *CreateTransactionRequest) GetTransactionType(db *gorm.DB) *models.Category {
	if c.TransactionType == "" {
		return scopes.GetOrCreateTransactionType("Unknown", db)
	}
	return scopes.GetOrCreateTransactionType(c.TransactionType, db)
}

func (c *CreateTransactionRequest) GetDateFormat() string {
	if c.DateFormat == "" {
		return time.RFC3339
	}
	return c.DateFormat
}

func (c *CreateTransactionRequest) FormatDate() time.Time {
	date, err := time.Parse(time.DateOnly, c.Date)
	if err != nil {
		fmt.Printf("Error parsing date %s", date)
		fmt.Println(err)
		return time.Now()
	}

	return date
}
