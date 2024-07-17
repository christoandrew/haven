package scopes

import "gorm.io/gorm"

func TransactionsTotalThisWeekVsLastWeek(accountId int, db *gorm.DB) *gorm.DB {
	query := `SELECT
    			ROUND(SUM(CASE WHEN transactions.date BETWEEN DATE_ADD(CURDATE(), INTERVAL(1-DAYOFWEEK(CURDATE())) DAY) AND DATE_ADD(CURDATE(), INTERVAL(7-DAYOFWEEK(CURDATE())) DAY) THEN transactions.amount ELSE 0 END), 2) AS this_week,
    			ROUND(SUM(CASE WHEN transactions.date BETWEEN DATE_ADD(CURDATE(), INTERVAL(-6-DAYOFWEEK(CURDATE())) DAY) AND DATE_ADD(CURDATE(), INTERVAL(0-DAYOFWEEK(CURDATE())) DAY) THEN transactions.amount ELSE 0 END), 2) AS last_week
			  FROM transactions
			  WHERE account_id = ?;`

	return db.Raw(query, accountId)
}
