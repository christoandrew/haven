package serializers

import (
	"github.com/christo-andrew/haven/internal/api/responses"
	"github.com/christo-andrew/haven/internal/models"
)

type TransactionSerializer struct {
	Data interface{}
	many bool
}

func NewTransactionSerializer(data interface{}, many bool) *TransactionSerializer {
	return &TransactionSerializer{
		Data: data,
		many: many,
	}
}

func (ts TransactionSerializer) Serialize() interface{} {
	switch ts.Data.(type) {
	case []models.Transaction:
		return ts.serializeTransactions()
	case models.Transaction:
		return ts.serializeTransaction()
	default:
		return nil
	}
}

func (ts TransactionSerializer) serializeTransactions() []responses.TransactionResponse {
	var response []responses.TransactionResponse
	for _, tx := range ts.Data.([]models.Transaction) {
		response = append(response, ts.serializeSingleTransaction(tx))
	}
	return response
}

func (ts TransactionSerializer) serializeTransaction() responses.TransactionResponse {
	return ts.serializeSingleTransaction(ts.Data.(models.Transaction))
}

func (ts TransactionSerializer) serializeSingleTransaction(tx models.Transaction) responses.TransactionResponse {
	return responses.TransactionResponse{
		TransactionID:     tx.ID,
		Amount:            tx.Amount,
		Description:       tx.Description,
		Date:              tx.Date.Unix(),
		AccountID:         tx.AccountID,
		TransactionType:   tx.TransactionType.Name,
		Category:          tx.Category.Name,
		Currency:          tx.Currency,
		Reference:         tx.Reference,
		Payee:             tx.Payee,
		TransactionStatus: tx.TransactionStatus,
	}
}
