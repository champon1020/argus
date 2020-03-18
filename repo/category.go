package repo

import (
	"database/sql"

	"github.com/champon1020/argus/service"
)

// Id: primary key
// Name: category name
type Category struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (category *Category) InsertCategory(tx *sql.Tx) (err error) {
	cmd := "INSERT INTO categories (id, name) " +
		"SELECT ?,? " +
		"WHERE NOT EXISTS (" +
		"SELECT name FROM categories WHERE name=?)"

	if _, err = tx.Exec(cmd,
		category.Id,
		category.Name,
		category.Name,
	); err != nil {
		CmdError.SetErr(err).AppendTo(Errors)
	}
	return
}

func (category *Category) UpdateCategory(tx *sql.Tx) (err error) {
	cmd := "UPDATE categories " +
		"SET name=? " +
		"WHERE id=? "

	if _, err = tx.Exec(cmd, category.Name, category.Id); err != nil {
		CmdError.SetErr(err).AppendTo(Errors)
	}
	return
}

func (category *Category) DeleteCategory(tx *sql.Tx) (err error) {
	cmd := "DELETE FROM categories WHERE id=?"
	if _, err = tx.Exec(cmd, category.Id); err != nil {
		CmdError.SetErr(err).AppendTo(Errors)
	}
	return
}

// Get the number of articles where category_id is equal to object.
func (category *Category) FindArticleNumByCategoryId(db *sql.DB) (articleNum int, err error) {
	query := "SELECT COUNT(article_id) FROM article_category WHERE category_id=?"

	var rows *sql.Rows
	defer RowsClose(rows)
	if rows, err = db.Query(query, category.Id); err != nil || rows == nil {
		QueryError.
			SetErr(err).
			SetValues("query", query).
			SetValues("args", category.Id).
			AppendTo(Errors)
		return
	}

	for rows.Next() {
		if err := rows.Scan(&articleNum); err != nil {
			ScanError.SetErr(err).AppendTo(Errors)
			break
		}
	}
	return
}

func (category *Category) Exist(tx *sql.Tx, option *service.QueryOption) (categoryId string, err error) {
	args := (*option).Args
	argsQuery := service.GenArgsQuery(*option)
	query := "SELECT id FROM categories " + argsQuery

	var rows *sql.Rows
	defer RowsClose(rows)
	if rows, err = tx.Query(query, args...); err != nil || rows == nil {
		QueryError.
			SetErr(err).
			SetValues("query", query).
			SetValues("args", args).
			AppendTo(Errors)
		return
	}

	for rows.Next() {
		if err := rows.Scan(&categoryId); err != nil {
			ScanError.SetErr(err).AppendTo(Errors)
			break
		}
	}
	return
}

// This is article category struct which is used only response.
// Difference of normal category struct is that this has property of 'ArticleNum'.
// ArticleNum is the number of articles related to this category.
type CategoryResponse struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	ArticleNum int    `json:"articleNum"`
}

func FindCategory(db *sql.DB, option *service.QueryOption) (categories []CategoryResponse, err error) {
	args := (*option).Args
	argsQuery := service.GenArgsQuery(*option)
	query := "SELECT * FROM categories " + argsQuery

	var rows *sql.Rows
	defer RowsClose(rows)
	if rows, err = db.Query(query, args...); err != nil || rows == nil {
		QueryError.
			SetErr(err).
			SetValues("query", query).
			SetValues("args", args).
			AppendTo(Errors)
		return
	}

	var (
		c          Category
		articleNum int
	)
	for rows.Next() {
		if err := rows.Scan(&c.Id, &c.Name); err != nil {
			ScanError.SetErr(err).AppendTo(Errors)
			break
		}
		if articleNum, err = c.FindArticleNumByCategoryId(db); err != nil {
			break
		}
		categories = append(
			categories,
			CategoryResponse{Id: c.Id, Name: c.Name, ArticleNum: articleNum},
		)
	}
	return
}
