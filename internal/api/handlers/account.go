package handlers

import (
	"io"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/christo-andrew/haven/internal/api/requests"
	"github.com/christo-andrew/haven/internal/api/schemas"
	"github.com/christo-andrew/haven/internal/api/serializers"
	"github.com/christo-andrew/haven/internal/models"
	"github.com/christo-andrew/haven/pkg/database/scopes"
	"github.com/christo-andrew/haven/pkg/pagination"
	utils "github.com/christo-andrew/haven/pkg/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetAllAccountsHandler GetAccounts godoc
// @Summary Get all accounts
// @Description Retrieve all accounts
// @Produce json
// @Success 200 {array} responses.AccountResponse
// @Router /accounts [get]
// @Param group_by_account_type query boolean false "Group by account type" Enums(true, false)
// @Tags accounts
// @Security AuthToken
// @Param Authorization header string true "Authorization"
func GetAllAccountsHandler(c *gin.Context, db *gorm.DB) {
	groupByAccountType, _ := strconv.ParseBool(c.Query("group_by_account_type"))
	accounts := getAccounts(groupByAccountType, db)
	response, err := serializers.NewAccountSerializer(accounts, true).Serialize()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func getAccounts(groupByAccountType bool, db *gorm.DB) interface{} {
	var accounts []models.Account
	if groupByAccountType {
		db.Model(&models.Account{}).Find(&accounts)
		result := make(map[string][]models.Account)
		for _, account := range accounts {
			result[account.BaseAccountType] = append(result[account.BaseAccountType], account)
		}
		return result
	}
	db.Find(&accounts)
	return accounts
}

// GetAccountHandler GetAccount godoc
// @Summary Get an account
// @Description Retrieve an account
// @Produce json
// @Param id path string true "Account ID"
// @Success 200 {object} responses.AccountResponse
// @Router /accounts/{id} [get]
// @Tags accounts
// @Security AuthToken
func GetAccountHandler(c *gin.Context, db *gorm.DB) {
	accountId, _ := strconv.Atoi(c.Param("id"))
	account := getAccount(accountId, db)
	serializer := serializers.NewAccountSerializer(account, false)
	response, err := serializer.Serialize()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}

func getAccount(accountID int, db *gorm.DB) models.Account {
	var account models.Account
	db.First(&account, accountID)
	return account
}

// GetAccountTransactionsHandler GetAccountTransactions godoc
// @Summary Get an account's transactions
// @Description Retrieve an account's transactions
// @Produce json
// @Param id path string true "Account ID"
// @Param page query int false "Page number"
// @Success 200 {array} responses.TransactionResponse
// @Router /accounts/{id}/transactions [get]
// @Tags accounts
// @Security AuthToken
func GetAccountTransactionsHandler(c *gin.Context, db *gorm.DB) {
	accountId, _ := strconv.Atoi(c.Param("id"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	transactions := scopes.GetTransactionsByAccountId(accountId, db)
	paginator := pagination.Pagination{Page: page, Limit: 20}
	var results []models.Transaction
	paginator.Paginate(transactions, models.Transaction{}).Find(&results)
	serializer := serializers.NewTransactionSerializer(results, true)

	response := pagination.Response{
		Results:    serializer.Serialize(),
		NextPage:   paginator.NextPage(),
		PrevPage:   paginator.PrevPage(),
		TotalCount: paginator.TotalCount,
		Limit:      paginator.Limit,
		Page:       paginator.Page,
		LastPage:   paginator.LastPage(),
	}
	c.JSON(http.StatusOK, response)
}

// CreateAccountHandler CreateAccount godoc
// @Summary Create an account
// @Description Create an account
// @Accept json
// @Produce json
// @Param account body requests.CreateBankAccountRequest true "Create Account Request"
// @Success 201 {object} responses.AccountResponse
// @Failure 400 {object} responses.ErrorResponse
// @Router /accounts/create [post]
// @Tags accounts
// @Security AuthToken
func CreateAccountHandler(c *gin.Context, db *gorm.DB) {
	var genericCreateAccountRequest requests.GenericCreateAccountRequest
	err := c.ShouldBindJSON(&genericCreateAccountRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error(err)})
		return
	}
	accountRequest, _ := requests.GetAccountRequest(&genericCreateAccountRequest)
	account := accountRequest.Account()
	account, err = createAccount(account, db)
	response, _ := serializers.NewAccountSerializer(account, false).Serialize()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, response)
}

func createAccount(account models.IAccount, db *gorm.DB) (models.IAccount, error) {
	result := db.Create(account)
	if result.Error != nil {
		return account, result.Error
	}
	return account, nil
}

// GetRecentTransactionsHandler Get RecentAccountTransactions godoc
// @Summary Get 5 recent transactions for an account
// @Description Get 5 recent transactions for an account
// @Accept json
// @Produce json
// @Param id path int true "Account ID"
// @Success 200 {array} responses.TransactionResponse
// @Router /accounts/{id}/transactions/recent [get]
// @Failure 400 {object} responses.ErrorResponse
// @Tags accounts
// @Security AuthToken
// @Param Authorization header string true "Authorization"
func GetRecentTransactionsHandler(c *gin.Context, db *gorm.DB) {
	accountId, _ := strconv.Atoi(c.Param("id"))
	var transactions []models.Transaction
	scopes.GetRecentTransactions(db, accountId, 4).Find(&transactions)
	response := serializers.NewTransactionSerializer(transactions, true).Serialize()
	c.JSON(http.StatusOK, response)
}

// UploadAccountTransactionsHandler Post Upload Account Transactions godoc
// @Summary Upload account transactions
// @Description Upload account transactions
// @Param id path int true "Account ID"
// @Param file formData file true "Transactions File"
// @Param transaction_schema formData string true "Transaction Schema"
// @Accept multipart/form-data
// @Produce json
// @Success 200 {array} responses.TransactionResponse
// @Router /accounts/{id}/transactions/upload [post]
// @Failure 400 {object} responses.ErrorResponse
// @Consumes multipart/form-data
// @Tags accounts
// @Security AuthToken
// @Param Authorization header string true "Authorization"
func UploadAccountTransactionsHandler(c *gin.Context, db *gorm.DB) {
	accountId, _ := strconv.Atoi(c.Param("id"))
	transactionSchemaType := c.PostForm("transaction_schema")
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	account := getAccount(accountId, db)
	if account.AccountName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Account not found"})
		return
	}
	transactionSchema := schemas.GetTransactionSchemaFromName(transactionSchemaType, &account)
	if transactionSchema == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Transaction schema not found"})
		return
	}
	openFile, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer func(openFile multipart.File) {
		err := openFile.Close()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	}(openFile)
	transactions := parseTransactionsFile(openFile, transactionSchema)
	db.Create(&transactions)
	c.JSON(http.StatusOK, transactions)
}

func parseTransactionsFile(file multipart.File, transactionSchema schemas.ITransactionSchema) []*models.Transaction {
	var transactions []*models.Transaction
	reader := io.Reader(file)
	transactionFromFile := utils.CSVToMap(reader)
	for _, transaction := range transactionFromFile {
		transaction := transactionSchema.Transaction(transaction)
		transactions = append(transactions, transaction)
	}

	return transactions
}
