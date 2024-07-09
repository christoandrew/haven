package handlers

import (
	"sort"
	"strconv"
	"time"

	"github.com/christo-andrew/haven/internal/models"
	"github.com/christo-andrew/haven/pkg/database/scopes"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type transactionsByYearAndMonth struct {
	Year            string  `json:"year"`
	Month           string  `json:"month"`
	Amount          float64 `json:"amount"`
	TransactionType string  `json:"transaction_type"`
}

type transactionsByCategory struct {
	Name   string  `json:"name"`
	Amount float64 `json:"amount"`
}

// TransactionsHistogramHandler TransactionHistogramData godoc
// @Summary Get transaction histogram data
// @Description Get transaction histogram data
// @ID get-transaction-histogram-data
// @Produce json
// @Param account_id path int true "Account ID"
// @Success 200 {array} any
// @Router /data/{account_id}/transactions/histogram [get]
// @Tags data
// @Security AuthToken
// @Param Authorization header string true "Authorization"
func TransactionsHistogramHandler(c *gin.Context, db *gorm.DB) {
	accountId, _ := strconv.Atoi(c.Param("account_id"))
	c.JSON(200, buildTransactionsHistogramData(accountId, db))
}

func buildTransactionsHistogramData(accountID int, db *gorm.DB) map[string]interface{} {
	var (
		result           []transactionsByYearAndMonth
		months           []int
		transactionTypes = hashset.New()
	)

	currentYear := 2024
	scopes.AccountTransactionsByYearAndMonth(accountID, 2024, db).Scan(&result)

	groupedByYearAndMonth := groupByYearAndMonth(result)
	forThisYear := groupedByYearAndMonth[strconv.Itoa(currentYear)]

	for monthStr := range forThisYear {
		monthInt, _ := strconv.Atoi(monthStr)
		months = append(months, monthInt)
	}

	// Sort the months
	sort.Ints(months)

	// Build data directly while iterating over months
	data := make([]map[string]interface{}, 0, len(months))
	for _, monthInt := range months {
		monthStr := strconv.Itoa(monthInt)
		transactionData := map[string]interface{}{
			"month":       time.Month(monthInt).String(),
			"monthNumber": monthInt,
		}

		for transactionType, amount := range forThisYear[monthStr] {
			transactionTypes.Add(transactionType)
			transactionData[transactionType] = amount
		}

		data = append(data, transactionData)
	}

	return map[string]interface{}{
		"data": data,
		"meta": map[string]interface{}{
			"transaction_types": transactionTypes.Values(),
			"colors":            models.TransactionTypeColors(),
		},
	}
}

func groupByYearAndMonth(transactions []transactionsByYearAndMonth) map[string]map[string]map[string]float64 {
	var transactionsMap map[string]map[string]map[string]float64
	transactionsMap = make(map[string]map[string]map[string]float64)

	for _, transaction := range transactions {
		if _, ok := transactionsMap[transaction.Year]; !ok {
			transactionsMap[transaction.Year] = make(map[string]map[string]float64)
		}

		if _, ok := transactionsMap[transaction.Year][transaction.Month]; !ok {
			transactionsMap[transaction.Year][transaction.Month] = make(map[string]float64)
		}

		transactionsMap[transaction.Year][transaction.Month][transaction.TransactionType] = transaction.Amount
	}

	return transactionsMap
}

// TransactionsSummaryHandler TransactionsSummaryData godoc
// @Summary Get transactions summary data
// @Description Get transactions summary data
//
//	Available filters: transaction_category
//	Available intervals: month, year, week
//	Available group_by: transaction_type
//	Available sort_by: amount
//	Available sort_order: asc, desc
//	Available limit: 10
//	Available offset: 0
//	Available account_id: 1
//
// @Param filter query string false "Filter", Enums(transaction_category, transaction_type, month, year, week)
// @Param interval query string false "Interval", Enums(month, year, week)
// @Param group_by query string false "Group by"
// @Param sort_by query string false "Sort by"
// @Param sort_order query string false "Sort order", Enums(asc, desc)
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Param account_id path int true "Account ID"
// @Tags data
// @Produce json
// @Success 200 {array} transactionSummaryData
// @Router /data/{account_id}/transactions/summary [get]
// @Security AuthToken
// @Param Authorization header string true "Authorization"
func TransactionsSummaryHandler(c *gin.Context, db *gorm.DB) {
	filter := c.Query("filter")
	switch filter {
	case "transaction_category":
		transactionsSummaryByTransactionCategoryHandler(c, db)
	}
}

func transactionsSummaryByTransactionCategoryHandler(c *gin.Context, db *gorm.DB) {
	//interval := c.Query("interval")
	accountId, _ := strconv.Atoi(c.Param("account_id"))
	data, meta := buildTransactionsSummaryByTransactionCategory(accountId, db)
	c.JSON(200, gin.H{
		"data": data,
		"meta": meta,
	})
}

func buildTransactionsSummaryByTransactionCategory(accountId int, db *gorm.DB) ([]transactionsByCategory, map[string]interface{}) {
	var result []transactionsByCategory
	scopes.GroupAccountTransactionsByTransactionCategory(accountId, db).Scan(&result)
	meta := map[string]interface{}{
		"colors": models.TransactionTypeColors(),
	}
	return result, meta
}
