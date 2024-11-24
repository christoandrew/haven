package models

import (
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
	"time"
)

// Accounts
type IAccount interface {
}

type Account struct {
	gorm.Model
	ID              int    `json:"id"`
	AccountName     string `json:"name"`
	AccountType     string `json:"account_type"`
	Currency        string `json:"currency"`
	UserID          uint   `json:"user_id"`
	Balance         float64
	BaseAccountType string        `json:"base_account_type"`
	BaseAccountID   int           `json:"base_account_id"`
	Transactions    []Transaction `gorm:"foreignKey:AccountID"`
}

type BankAccount struct {
	gorm.Model
	Account Account `gorm:"polymorphic:BaseAccount;"`
}

type CreditCardAccount struct {
	gorm.Model
	Account Account `gorm:"polymorphic:BaseAccount;"`
}

// Budget
type Budget struct {
	gorm.Model
	Id          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	UserId      uint      `json:"user_id"`
	User        User      `json:"user"`
	Category    Category  `json:"category"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	CategoryID  uint      `json:"category_id"`
}

type BudgetCategory struct {
	gorm.Model
	BudgetID   int       `json:"budget_id"`
	Budget     *Budget   `json:"budget"`
	CategoryID int       `json:"category_id"`
	Category   *Category `json:"category"`
}

type Category struct {
	gorm.Model
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Context     string `json:"context"`
	ContextType string `json:"context_type"`
}

type RealEstateAccount struct {
	gorm.Model
	Account Account `gorm:"polymorphic:BaseAccount;"`
}

type Tag struct {
	gorm.Model
	ID   int    `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
}

type Transaction struct {
	gorm.Model
	ID                int       `json:"id"`
	Amount            float64   `json:"amount"`
	Currency          string    `json:"currency"`
	Payee             string    `json:"payee"`
	Reference         string    `json:"reference"`
	Date              time.Time `json:"date"`
	Description       string    `json:"description"`
	AccountID         int       `json:"account_id"`
	Account           Account   `gorm:"foreignKey:AccountID"`
	CategoryID        int       `json:"category_id"`
	Category          Category  `gorm:"foreignKey:CategoryID"`
	TransactionTypeID int       `json:"transaction_type_id"`
	TransactionType   Category  `gorm:"foreignKey:TransactionTypeID"`
	TransactionStatus string    `json:"transaction_status"`
	Tags              []Tag     `gorm:"many2many:transaction_tags;"`
}

func TransactionTypeColors() map[string]string {
	return map[string]string{
		"debit":    "#FDA403",
		"credit":   "#898121",
		"deposit":  "#E5C287",
		"withdraw": "#E8751A",
	}
}

type User struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password" gorm:"type:varchar(256)"`
}

func (user *User) GenerateToken() (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"id":    user.ID,
	})

	return claims.SignedString([]byte("secret"))
}

func (user *User) GetFullName() string {
	return user.FirstName + " " + user.LastName
}