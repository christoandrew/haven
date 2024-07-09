package responses

import "github.com/christo-andrew/haven/internal/models"

type ErrorResponse struct {
	Message string `json:"message"`
}

type CreateUserResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func (createUserResponse CreateUserResponse) FromUser(user *models.User) CreateUserResponse {
	createUserResponse.ID = user.ID
	createUserResponse.Email = user.Email
	createUserResponse.Name = user.GetFullName()
	return createUserResponse
}

type UserResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type AccountResponse struct {
	ID          int     `json:"id"`
	AccountName string  `json:"name"`
	Currency    string  `json:"currency"`
	Balance     float64 `json:"balance"`
	AccountType string  `json:"account_type"`
	Category    string  `json:"category"`
}

type CategoryResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Context     string `json:"context"`
	ContextType string `json:"context_type"`
}

type TransactionResponse struct {
	TransactionID     int     `json:"id"`
	Amount            float64 `json:"amount"`
	Currency          string  `json:"currency"`
	Date              int64   `json:"date"`
	Description       string  `json:"description"`
	AccountID         int     `json:"account_id"`
	TransactionType   string  `json:"transaction_type"`
	Category          string  `json:"category"`
	TransactionStatus string  `json:"transaction_status"`
	Reference         string  `json:"reference"`
	Payee             string  `json:"payee"`
}

type TagResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type TransactionSchema struct {
	Name       string `yaml:"name" json:"name"`
	DateFormat string `yaml:"date_format" json:"date_format"`
	Mapping    []struct {
		Name    string      `yaml:"name" json:"name"`
		Type    string      `yaml:"type" json:"type"`
		Column  string      `yaml:"column" json:"column"`
		Default interface{} `yaml:"default" json:"default"`
	} `yaml:"mapping" json:"mapping"`
	Computations []struct {
		Name    string `yaml:"name" json:"name"`
		Formula string `yaml:"formula" json:"formula"`
	}
}
