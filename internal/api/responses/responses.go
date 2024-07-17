package responses

import (
	"github.com/christo-andrew/haven/internal/models"
)

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

type BudgetResponse struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
	StartDate   string  `json:"start_date"`
	EndDate     string  `json:"end_date"`
	Category    string  `json:"category"`
}

func (budgetResponse BudgetResponse) FromBudget(budget models.Budget) *BudgetResponse {
	budgetResponse.Id = int(budget.Id)
	budgetResponse.Name = budget.Name
	budgetResponse.Description = budget.Description
	budgetResponse.Amount = budget.Amount
	budgetResponse.StartDate = budget.StartDate.Format("2006-01-02")
	budgetResponse.EndDate = budget.EndDate.Format("2006-01-02")
	budgetResponse.Category = budget.Category.Name
	return &budgetResponse
}

type PercentageOfTotalAmountByTransactionResponse struct {
	Category   string  `json:"category"`
	Amount     float64 `json:"amount"`
	Percentage float64 `json:"percentage"`
}

// AccountStatisticsResponse holds the overall statistics for an account.
type AccountStatisticsResponse struct {
	TotalBalance float64               `json:"total_balance"`
	TotalIncome  float64               `json:"total_income"`
	TotalExpense float64               `json:"total_expense"`
	Transactions TransactionStatistics `json:"transactions"`
}

// TransactionStatistics holds detailed statistics about transactions.
type TransactionStatistics struct {
	ThisWeekVsLastWeek WeekComparison `json:"this_week_vs_last_week"`
}

// WeekComparison holds the comparison data between this week and last week.
type WeekComparison struct {
	ThisWeek         float64 `json:"this_week"`
	LastWeek         float64 `json:"last_week"`
	Change           float64 `json:"change"`
	PercentageChange float64 `json:"percentage_change"`
}

func (weekComparison *WeekComparison) CalculateChange() {
	weekComparison.Change = weekComparison.ThisWeek - weekComparison.LastWeek
}
