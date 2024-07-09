package serializers

import (
	"github.com/christo-andrew/haven/internal/api/responses"
	"github.com/christo-andrew/haven/internal/models"
	"github.com/christo-andrew/haven/pkg/errors"
)

type AccountSerializer struct {
	Data interface{}
	many bool
}

func NewAccountSerializer(data interface{}, many bool) *AccountSerializer {
	return &AccountSerializer{
		Data: data,
		many: many,
	}
}

func (as AccountSerializer) Serialize() (interface{}, error) {
	switch as.Data.(type) {
	case map[string][]models.Account:
		return as.serializeGrouped()
	case []models.Account:
		return as.serializeMany(as.Data)
	case models.Account:
		return as.serializeSingle(as.Data)
	default:
		return nil, errors.InvalidDataError()
	}
}

func (as AccountSerializer) serializeMany(obj interface{}) (interface{}, error) {
	accounts, ok := obj.([]models.Account)
	if !ok {
		return nil, errors.InvalidDataError()
	}

	var response []*responses.AccountResponse
	for _, account := range accounts {
		data, err := as.serializeSingle(account)
		if err != nil {
			return nil, err
		}
		response = append(response, data)
	}
	return response, nil
}

func (as AccountSerializer) serializeSingle(obj interface{}) (*responses.AccountResponse, error) {
	if obj == nil {
		return nil, errors.InvalidDataError()
	}
	account, ok := obj.(models.Account)
	if !ok {
		return nil, errors.InvalidDataError()
	}

	return &responses.AccountResponse{
		ID:          account.ID,
		AccountName: account.AccountName,
		Currency:    account.Currency,
		Balance:     account.Balance,
		AccountType: account.AccountType,
		Category:    account.BaseAccountType,
	}, nil
}

func (as AccountSerializer) serializeGrouped() (interface{}, error) {
	if as.Data == nil {
		return nil, errors.InvalidDataError()
	}

	accounts, ok := as.Data.(map[string][]models.Account)
	if !ok {
		return nil, errors.InvalidDataError()
	}

	accountMap := make(map[string][]*responses.AccountResponse)
	for key, value := range accounts {
		var response []*responses.AccountResponse
		for _, account := range value {
			data, err := as.serializeSingle(account)
			if err != nil {
				return nil, err
			}
			response = append(response, data)
		}
		accountMap[key] = response
	}
	return accountMap, nil
}
