package handlers

import (
	"github.com/christo-andrew/haven/internal/api/requests"
	"github.com/christo-andrew/haven/internal/api/responses"
	"github.com/christo-andrew/haven/internal/api/serializers"
	"github.com/christo-andrew/haven/internal/models"
	"github.com/christo-andrew/haven/pkg/database/scopes"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
	"net/http"
	"os"
	"strconv"
)

// GetAllTransactionsHandler GetAllTransactions godoc
// @Summary Get all transactions
// @Description Retrieve all transactions
// @Produce json
// @Success 200 {array} responses.TransactionResponse
// @Router /transactions [get]
// @Tags transactions
// @Security AuthToken
// @Param Authorization header string true "Authorization"
func GetAllTransactionsHandler(c *gin.Context, db *gorm.DB) {
	serializer := serializers.NewTransactionSerializer(getAllTransactions(db), true)
	c.JSON(http.StatusOK, serializer.Serialize())
}

func getAllTransactions(db *gorm.DB) []models.Transaction {
	var transactions []models.Transaction
	scopes.GetAllTransactions(db).Find(&transactions)
	return transactions
}

// GetTransactionHandler GetTransaction godoc
// @Summary Get a transaction
// @Description Retrieve a transaction
// @Produce json
// @Param id path string true "Transaction ID"
// @Success 200 {object} responses.TransactionResponse
// @Router /transactions/{id} [get]
// @Tags transactions
// @Security AuthToken
// @Param Authorization header string true "Authorization"
func GetTransactionHandler(c *gin.Context, db *gorm.DB) {
	transactionId, _ := strconv.Atoi(c.Param("id"))
	response := serializers.NewTransactionSerializer(getTransaction(transactionId, db), false).Serialize()
	c.JSON(http.StatusOK, response)
}

func getTransaction(transactionID int, db *gorm.DB) models.Transaction {
	var transaction models.Transaction
	scopes.GetTransactionById(transactionID, db).Find(&transaction)
	return transaction
}

// CreateAccountTransactionHandler CreateTransaction godoc
// @Summary Create a transaction
// @Description Create a transaction
// @Accept json
// @Produce json
// @Param transaction body requests.CreateTransactionRequest true "Create Transaction Request"
// @Success 201 {object} responses.TransactionResponse
// @Failure 400 {object} responses.ErrorResponse
// @Router /transactions/create [post]
// @Tags transactions
// @Security AuthToken
// @Param Authorization header string true "Authorization"
func CreateAccountTransactionHandler(c *gin.Context, db *gorm.DB) {
	batchCreate := c.Query("batch_create")
	if batchCreate == "true" {
		createBatchTransactions(c, db)
		return
	}
	var transactionRequest requests.CreateTransactionRequest
	err := c.ShouldBindJSON(&transactionRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	transaction, err := createTransaction(&transactionRequest, db)
	response := serializers.NewTransactionSerializer(transaction, false).Serialize()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, response)
}

// AddTransactionTagHandler AddTransactionTag godoc
// @Summary Add a tag to a transaction
// @Description Add a tag to a transaction
// @Accept json
// @Produce json
// @Param id path string true "Transaction ID"
// @Param tag body requests.CreateTagRequest true "Create Tag Request"
// @Success 201 {object} responses.TagResponse
// @Failure 400 {object} responses.ErrorResponse
// @Router /transactions/{id}/tags [post]
// @Tags transactions
// @Security AuthToken
// @Param Authorization header string true "Authorization"
func AddTransactionTagHandler(c *gin.Context, db *gorm.DB) {
	transactionId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var createTagRequest requests.CreateTagRequest
	err = c.ShouldBindJSON(&createTagRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	transaction := getTransaction(transactionId, db)
	tag := scopes.GetOrCreateTransactionTag(createTagRequest.Name, db)
	response := serializers.NewTagSerializer(*tag, false).Serialize()
	err = db.Model(&transaction).Association("Tags").Append(tag)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, response)
}

// GetTransactionSchemasHandler GetTransactionSchemas godoc
// @Summary Get transaction schemas
// @Description Retrieve transaction schemas
// @Produce json
// @Success 200 {array} responses.TransactionSchema
// @Router /transactions/schemas [get]
// @Tags transactions
// @Security AuthToken
// @Param Authorization header string true "Authorization"
func GetTransactionSchemasHandler(c *gin.Context) {
	schemaFile := "internal/api/schemas/transaction.yml"
	data, err := os.ReadFile(schemaFile)
	var transactionSchema []responses.TransactionSchema
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = yaml.Unmarshal(data, &transactionSchema)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transactionSchema)
}

// GetTransactionTagsHandler GetTransactionTags godoc
// @Summary Get all tags for a transaction
// @Description Retrieve all tags for a transaction
// @Produce json
// @Param id path string true "Transaction ID"
// @Success 200 {array} responses.TagResponse
// @Router /transactions/{id}/tags [get]
// @Tags transactions
// @Security AuthToken
// @Param Authorization header string true "Authorization"
func GetTransactionTagsHandler(c *gin.Context, db *gorm.DB) {
	var transaction models.Transaction
	transactionId, _ := strconv.Atoi(c.Param("id"))
	scopes.GetTransactionByIdWithTags(transactionId, db).Find(&transaction)
	tags := transaction.Tags
	response := serializers.NewTagSerializer(tags, true).Serialize()
	c.JSON(http.StatusOK, response)
}

func createTransaction(transactionRequest *requests.CreateTransactionRequest, db *gorm.DB) (*models.Transaction, error) {
	transaction := transactionRequest.Transaction(db)
	transaction.Category = *transactionRequest.GetCategory(db)
	transaction.TransactionType = *transactionRequest.GetTransactionType(db)
	result := db.Create(transaction)

	return transaction, result.Error
}

func createBatchTransactions(c *gin.Context, db *gorm.DB) {
	var transactionRequests []requests.CreateTransactionRequest
	err := c.ShouldBindJSON(&transactionRequests)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	transactions, err := createTransactions(transactionRequests, db)
	response := serializers.NewTransactionSerializer(transactions, true).Serialize()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, response)
}

func createTransactions(transactionRequests []requests.CreateTransactionRequest, db *gorm.DB) ([]models.Transaction, error) {
	var transactions []models.Transaction
	for _, transactionRequest := range transactionRequests {
		transaction := transactionRequest.Transaction(db)
		transaction.Category = *transactionRequest.GetCategory(db)
		transaction.TransactionType = *transactionRequest.GetTransactionType(db)
		result := db.Create(transaction)
		if result.Error != nil {
			return transactions, result.Error
		}
		transactions = append(transactions, *transaction)
	}
	return transactions, nil
}
