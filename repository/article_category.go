package repository

import "database/sql"

func FindArticleByCategoryId(db *sql.DB, categoryNames []string) (articles []Article, err error) {
	query := "SELECT * FROM articles " +
		"WHERE id IN (" +
		"SELECT article_id FROM article_category " +
		"WHERE category_id IN (" +
		"SELECT id FROM categories " +
		"WHERE "

	var args []interface{}
	for i, cn := range categoryNames {
		if i != 0 {
			query += "AND "
		}
		query += "name=? "
		args = append(args, cn)
	}
	query += "))"

	rows, err := db.Query(query, args...)
	defer func() {
		if rows == nil {
			return
		}
		if err := rows.Close(); err != nil {
			logger.ErrorPrintf(err)
		}
	}()

	if err != nil {
		logger.ErrorMsgPrintf("Unable to scan rows because rows is nil", err)
		return
	}

	if rows == nil {
		logger.ErrorMsgPrintf("Unable to scan rows because rows is nil", err)
		return
	}

	for rows.Next() {
		var a Article
		if err := rows.Scan(
			&a.Id,
			&a.Title,
			&a.CreateDate,
			&a.UpdateDate,
			&a.ContentUrl,
			&a.ImageUrl,
			&a.Private); err != nil {
			logger.ErrorPrintf(err)
		}
		if a.Categories, err = a.FindCategoryByArticleId(db); err != nil {
			logger.ErrorPrintf(err)
			break
		}
		articles = append(articles, a)
	}
	return
}
