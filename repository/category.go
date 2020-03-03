package repository

import "database/sql"

type Category struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func (category *Category) InsertCategory(tx *sql.Tx) (err error) {
	cmd := "INSERT INTO categories (id, name) " +
		"SELECT ?, ? " +
		"WHERE NOT EXISTS (" +
		"SELECT name FROM categories WHERE name=?)"

	_, err = tx.Exec(cmd,
		category.Id,
		category.Name,
		category.Name)

	if err != nil {
		logger.ErrorMsgPrintf("InsertCategory", err)
	}
	return
}

func (category *Category) UpdateCategory(tx *sql.Tx) (err error) {
	cmd := "UPDATE categories " +
		"SET name=? " +
		"WHERE id=? "

	_, err = tx.Exec(cmd,
		category.Name,
		category.Id)

	if err != nil {
		logger.ErrorPrintf(err)
	}
	return
}

func (category *Category) DeleteCategory(tx *sql.Tx) (err error) {
	cmd := "DELETE FROM categories WHERE id=?"
	_, err = tx.Exec(cmd, category.Id)
	if err != nil {
		logger.ErrorPrintf(err)
	}
	return
}

func (category *Category) DeleteArticleCategoryByCategory(tx *sql.Tx) (err error) {
	cmd := "DELETE FROM article_category WHERE category_id=?"
	_, err = tx.Exec(cmd, category.Id)
	if err != nil {
		logger.ErrorPrintf(err)
	}
	return
}

func (category *Category) FindArticleNumFromArticleCategory(db *sql.DB) (articleNum int, err error) {
	cmd := "SELECT COUNT(article_id) FROM article_category WHERE category_id=?"

	rows, err := db.Query(cmd, category.Id)
	defer func() {
		if rows == nil {
			return
		}
		if err := rows.Close(); err != nil {
			logger.ErrorPrintf(err)
		}
	}()

	if err != nil {
		logger.ErrorPrintf(err)
		return
	}

	if rows == nil {
		logger.ErrorMsgPrintf("Unable to scan rows because rows is nil", err)
		return
	}

	for rows.Next() {
		if err := rows.Scan(
			&articleNum); err != nil {
			logger.ErrorPrintf(err)
		}
	}
	return
}

type CategoryResponse struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	ArticleNum int    `json:"articleNum"`
}

func (category *Category) FindCategory(db *sql.DB, argsFlg uint32) (categories []CategoryResponse, err error) {
	args := GenArgsSlice(argsFlg, category)
	whereQuery := GenArgsQuery(argsFlg, category)
	query := "SELECT * FROM categories " + whereQuery + "ORDER BY id"

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
		logger.ErrorPrintf(err)
		return
	}

	if rows == nil {
		logger.ErrorMsgPrintf("Unable to scan rows because rows is nil", err)
		return
	}

	for rows.Next() {
		var (
			c          Category
			articleNum int
		)
		if err := rows.Scan(
			&c.Id,
			&c.Name); err != nil {
			logger.ErrorPrintf(err)
		}
		if articleNum, err = c.FindArticleNumFromArticleCategory(db); err != nil {
			return
		}

		categories = append(categories, CategoryResponse{Id: c.Id, Name: c.Name, ArticleNum: articleNum})
	}
	return
}
