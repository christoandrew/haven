package requests

import (
	"errors"

	"github.com/christo-andrew/haven/internal/models"
)

type GenericCreateAccountRequest struct {
	AccountName string  `json:"account_name"`
	AccountType string  `json:"account_type"`
	Currency    string  `json:"currency"`
	UserID      uint    `json:"user_id"`
	Balance     float64 `json:"balance"`
	Category    string  `json:"category"`
}

type CreateBankAccountRequest struct {
	*GenericCreateAccountRequest
}

type CreateCreditCardAccountRequest struct {
	*GenericCreateAccountRequest
}

type CreateRealEstateAccountRequest struct {
	*GenericCreateAccountRequest
}

type CreateLoanAccountRequest struct {
	*GenericCreateAccountRequest
}

type CreateInvestmentAccountRequest struct {
	*GenericCreateAccountRequest
}

type CreateAssetAccountRequest struct {
	*GenericCreateAccountRequest
}

type CreateIncomeAccountRequest struct {
	*GenericCreateAccountRequest
}

type CreateLiabilityAccountRequest struct {
	*GenericCreateAccountRequest
}

type CreateExpensesAccountRequest struct {
	*GenericCreateAccountRequest
}

// CreateAccountRequest defines the interface for creating different types of accounts
type CreateAccountRequest interface {
	Account() models.IAccount
}

func (c *GenericCreateAccountRequest) Account() models.IAccount {
	return nil // Return nil or handle error as per your logic
}

func GetAccountRequest(genericReq *GenericCreateAccountRequest) (CreateAccountRequest, error) {
	switch genericReq.Category {
	case "bank":
		return &CreateBankAccountRequest{GenericCreateAccountRequest: genericReq}, nil
	case "credit_card":
		return &CreateCreditCardAccountRequest{GenericCreateAccountRequest: genericReq}, nil
	case "real_estate":
		return &CreateRealEstateAccountRequest{GenericCreateAccountRequest: genericReq}, nil
	case "loan":
		return &CreateRealEstateAccountRequest{GenericCreateAccountRequest: genericReq}, nil
	case "investment":
		return &CreateRealEstateAccountRequest{GenericCreateAccountRequest: genericReq}, nil
	case "asset":
		return &CreateRealEstateAccountRequest{GenericCreateAccountRequest: genericReq}, nil
	case "income":
		return &CreateRealEstateAccountRequest{GenericCreateAccountRequest: genericReq}, nil
	case "liability":
		return &CreateRealEstateAccountRequest{GenericCreateAccountRequest: genericReq}, nil
	case "expenses":
		return &CreateRealEstateAccountRequest{GenericCreateAccountRequest: genericReq}, nil
	}
	return nil, errors.New("invalid account type")
}

func (c CreateBankAccountRequest) Account() models.IAccount {
	return &models.BankAccount{
		Account: c.createAccount(),
	}
}

func (c CreateCreditCardAccountRequest) Account() models.IAccount {
	return &models.CreditCardAccount{
		Account: c.createAccount(),
	}
}

func (c CreateRealEstateAccountRequest) Account() models.IAccount {
	return &models.RealEstateAccount{
		Account: c.createAccount(),
	}
}

// createAccount creates an account based on the request
func (c *GenericCreateAccountRequest) createAccount() models.Account {
	return models.Account{
		AccountName: c.AccountName,
		AccountType: c.AccountType,
		Currency:    c.Currency,
		UserID:      c.UserID,
		Balance:     c.Balance,
	}
}
