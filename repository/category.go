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

func (category *Category) FindCategory(db *sql.DB, argsFlg uint32) (categories []Category, err error) {
	args := GenArgsSlice(argsFlg, category)
	whereQuery := GenArgsQuery(argsFlg, category)
	query := "SELECT * FROM categories " + whereQuery + "ORDER BY id"

	rows, err := db.Query(query, args...)
	defer func() {
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
		var c Category
		if err := rows.Scan(
			&c.Id,
			&c.Name); err != nil {
			logger.ErrorPrintf(err)
		}
		categories = append(categories, c)
	}
	return
}
