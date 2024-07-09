package scopes

import "gorm.io/gorm"

func AccountTransactionsByYearAndMonth(accountId int, year int, db *gorm.DB) *gorm.DB {
	query := `SELECT
				MONTH(transactions.date) as month,
				YEAR(transactions.date) as year,
				SUM(transactions.amount) as amount,
				transaction_types.name as transaction_type
			  FROM transactions
			  INNER JOIN categories AS transaction_types ON transactions.transaction_type_id = transaction_types.id
			  INNER JOIN categories AS transaction_categories ON transactions.category_id = transaction_categories.id
			  WHERE transactions.account_id = ? AND YEAR(transactions.date) = ?
			  GROUP BY MONTH(transactions.date),YEAR(transactions.date), transaction_type
			  ORDER BY MONTH(transactions.date) ASC, YEAR(transactions.date) DESC;`

	return db.Raw(query, accountId, year)
}

func TotalAmountByTransactionTypeGroupedByYearAndMonth(accountId int, db *gorm.DB) *gorm.DB {
	query := `SELECT
				MONTH(transactions.date)
				YEAR(transactions.date)
       			SUM(ABS(transactions.amount)) AS amount
              FROM transactions
              INNER JOIN categories AS transaction_types ON transactions.transaction_type_id = transaction_types.id
              INNER JOIN categories AS transaction_categories ON transactions.category_id = transaction_categories.id
              WHERE transactions.account_id = ?
              ORDER BY month DESC
              GROUP BY YEAR(transactions.date), MONTH(transactions.date) ;`
	return db.Raw(query, accountId)
}

func GetTransactionById(id int, db *gorm.DB) *gorm.DB {
	return db.Preload("TransactionType").Preload("Category").Where("id = ?", id)
}

func GroupAccountTransactionsByTransactionCategory(accountId int, db *gorm.DB) *gorm.DB {
	query := `SELECT
				transaction_categories.name,
				SUM(transactions.amount) AS amount
			  FROM transactions
			  INNER JOIN categories AS transaction_types ON transactions.transaction_type_id = transaction_types.id
			  INNER JOIN categories AS transaction_categories ON transactions.category_id = transaction_categories.id
			  WHERE transactions.account_id = ?
			  GROUP BY transaction_categories.name;`

	return db.Raw(query, accountId)
}

func GetTransactionByIdWithTags(transactionId int, db *gorm.DB) *gorm.DB {
	return db.Preload("Tags").Where("id = ?", transactionId)
}

func GetAllTransactions(db *gorm.DB) *gorm.DB {
	return db.Preload("TransactionType").Preload("Category")
}

func GetTransactionsByAccountId(accountId int, db *gorm.DB) *gorm.DB {
	return db.Scopes(GetAllTransactions).Where("account_id = ?", accountId).Order("id DESC")
}

func GetTransactionsByCategoryId(categoryId int, db *gorm.DB) *gorm.DB {
	return db.Scopes(GetAllTransactions).Where("category_id = ?", categoryId)
}

func GetTransactionsByTransactionTypeId(transactionTypeId int, db *gorm.DB) *gorm.DB {
	return db.Scopes(GetAllTransactions).Where("transaction_type_id = ?", transactionTypeId)
}

func GetTransactionsByDateRange(startDate string, endDate string, db *gorm.DB) *gorm.DB {
	return db.Scopes(GetAllTransactions).Where("date BETWEEN ? AND ?", startDate, endDate)
}

func GetTransactionsByDateRangeAndAccountId(startDate string, endDate string, accountId int, db *gorm.DB) *gorm.DB {
	return GetTransactionsByAccountId(accountId, db).Where("date BETWEEN ? AND ?", startDate, endDate)
}

func GetTransactionsByDateRangeAndCategoryId(startDate string, endDate string, categoryId int, db *gorm.DB) *gorm.DB {
	return GetTransactionsByCategoryId(categoryId, db).Where("date BETWEEN ? AND ?", startDate, endDate)
}

func GetTransactionsByDateRangeAndTransactionTypeId(startDate string, endDate string, transactionTypeId int, db *gorm.DB) *gorm.DB {
	return GetTransactionsByTransactionTypeId(transactionTypeId, db).Where("date BETWEEN ? AND ?", startDate, endDate)
}

func GetTransactionsByDateRangeAndAccountIdAndCategoryId(startDate string, endDate string, accountId int, categoryId int, db *gorm.DB) *gorm.DB {
	return GetTransactionsByDateRangeAndAccountId(startDate, endDate, accountId, db).Where("category_id = ?", categoryId)
}
