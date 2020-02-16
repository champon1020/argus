package repository

import "database/sql"

func GetEmptyMinId(db *sql.DB, tableName string, numOfId int) (res []int, err error) {
	query := "SELECT (id+1) FROM " +
		tableName + " " +
		"WHERE (id+1) NOT IN " +
		"(SELECT id FROM " +
		tableName + " " +
		") LIMIT ?"

	rows, err := db.Query(query, numOfId)
	defer func() {
		if err = rows.Close(); err != nil {
			logger.ErrorPrintf(err)
			return
		}
	}()

	if err != nil {
		logger.ErrorPrintf(err)
		return
	}

	if rows == nil {
		logger.ErrorPrintf(err)
		return
	}

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			logger.ErrorPrintf(err)
			break
		}
		res = append(res, id)
	}
	return
}
