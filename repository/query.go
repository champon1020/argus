package repository

import "database/sql"

func GetEmptyMinId(db *sql.DB, tableName string, numOfId int) (res []int) {
	query := "SELECT (id+1) FROM ? " +
		"WHERE (id+1) NOT INT (" +
		"SELECT id FROM articles) LIMIT ?"

	rows, err := db.Query(query, tableName, numOfId)
	defer func() {
		if err := rows.Close(); err != nil {
			logger.ErrorPrintf(err)
		}
	}()

	if err != nil {
		logger.ErrorPrintf(err)
	}

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			logger.ErrorPrintf(err)
		}
		res = append(res, id)
	}
	return
}
