package handlers

import (
	"github.com/christo-andrew/haven/internal/api/responses"
	"github.com/christo-andrew/haven/pkg/pagination"
	"github.com/christo-andrew/haven/pkg/utils"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"

	"github.com/christo-andrew/haven/internal/api/requests"
	"github.com/christo-andrew/haven/internal/api/schemas"
	"github.com/christo-andrew/haven/internal/api/serializers"
	"github.com/christo-andrew/haven/internal/models"
	"github.com/christo-andrew/haven/pkg/database/scopes"
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
// @Param limit query int false "Limit"
// @Param from query string false "From" Format(YYYY-MM-DD)
// @Param to query string false "To" Format(YYYY-MM-DD)
// @Param unixTime query boolean false "Unix Time"
// @Success 200 {array} responses.TransactionResponse
// @Router /accounts/{id}/transactions [get]
// @Tags accounts
// @Security AuthToken
// @Param Authorization header string true "Authorization"
func GetAccountTransactionsHandler(c *gin.Context, db *gorm.DB) {
	accountId, _ := strconv.Atoi(c.Param("id"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	from := c.Query("from")
	to := c.Query("to")
	unixTime, _ := strconv.ParseBool(c.Query("unixTime"))
	paginator := pagination.Pagination{Page: page, Limit: limit}
	var results []models.Transaction
	var transactions *gorm.DB

	if from != "" && to != "" {
		var fromDate time.Time
		var toDate time.Time
		if unixTime {
			fromDate = utils.ConvertToUnixTime(from)
			toDate = utils.ConvertToUnixTime(to)
		}
		transactions = scopes.GetTransactionsByDateRangeAndAccountId(fromDate, toDate, accountId, db)
	} else {
		transactions = scopes.GetTransactionsByAccountId(accountId, db)
	}

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

// PercentageOfTotalAmountByTransactionHandler PercentageOfTotalAmountByTransaction godoc
// @Summary Get percentage of total amount by transaction category
// @Description Retrieve percentage of total amount by transaction category
// @Produce json
// @Param id path string true "Account ID"
// @Param limit query int false "Limit"
// @Param filter query string false "Filter" Enums(category)
// @Success 200 {array} responses.PercentageOfTotalAmountByTransactionResponse
// @Router /accounts/{id}/transactions/percentage [get]
// @Tags accounts
// @Security AuthToken
// @Param Authorization header string true "Authorization"
func PercentageOfTotalAmountByTransactionHandler(c *gin.Context, db *gorm.DB) {
	accountId, _ := strconv.Atoi(c.Param("id"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))
	filter := c.DefaultQuery("filter", "category")
	result := getPercentageOfTotalAmountBy(accountId, limit, db, filter)
	c.JSON(http.StatusOK, result)
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

// AccountStatisticsHandler Get Account Statistics godoc
// @Summary Get account statistics
// @Description Get account statistics
// @Accept json
// @Produce json
// @Param id path int true "Account ID"
// @Success 200 {object} responses.AccountStatisticsResponse
// @Router /accounts/{id}/statistics [get]
// @Failure 400 {object} responses.ErrorResponse
// @Tags accounts
// @Security AuthToken
// @Param Authorization header string true "Authorization"
func AccountStatisticsHandler(c *gin.Context, db *gorm.DB) {
	accountId, _ := strconv.Atoi(c.Param("id"))
	statistics := getStatistics(accountId, db)
	c.JSON(http.StatusOK, statistics)
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
	transactionSchema := schemas.GetTransactionSchemaFromName(transactionSchemaType, &account, db)
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

func getPercentageOfTotalAmountBy(accountId int, limit int, db *gorm.DB, filter string) interface{} {
	var result []*responses.PercentageOfTotalAmountByTransactionResponse
	switch filter {
	case "category":
		scopes.PercentageOfTotalAmountByTransactionCategory(accountId, limit, db).Scan(&result)
		return result
	}
	return nil
}

type TotalTransactionLastWeekVsThisWeek struct {
	LastWeekTotal    float64 `json:"last_week"`
	ThisWeek         float64 `json:"this_week"`
	PercentageChange float64 `json:"percentage_change"`
}

func getStatistics(accountId int, db *gorm.DB) interface{} {
	response := responses.AccountStatisticsResponse{}
	var weekComparison responses.WeekComparison
	scopes.TransactionsTotalThisWeekVsLastWeek(accountId, db).Scan(&weekComparison)
	weekComparison.CalculateChange()
	response.Transactions.ThisWeekVsLastWeek = weekComparison
	return response
}
