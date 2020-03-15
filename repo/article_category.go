package repo

import (
	"database/sql"

	"github.com/champon1020/argus/service"
)

func FindArticleByCategoryId(
	db *sql.DB,
	categoryNames []string,
	argsFlg uint32,
	ol OffsetLimit,
) (articles []Article, err error) {
	query := "SELECT * FROM articles " +
		"WHERE id IN (" +
		"SELECT article_id FROM article_category " +
		"WHERE category_id IN (" +
		"SELECT id FROM categories "

	args, whereQuery := GenArgsFromStrSlice(categoryNames)
	query += whereQuery
	query += ")) "

	args = append(args, service.GenArgsSlice(argsFlg, Article{}, ol)...)
	_, limitQuery := service.GenArgsQuery(argsFlg, Article{})
	query += limitQuery + "OFFSET ?"

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

	var a Article
	for rows.Next() {
		if err := rows.Scan(
			&a.Id,
			&a.Title,
			&a.CreateDate,
			&a.UpdateDate,
			&a.ContentHash,
			&a.ImageHash,
			&a.Private,
		); err != nil {
			ScanError.SetErr(err).AppendTo(Errors)
			break
		}
		if a.Categories, err = a.FindCategoryByArticleId(db); err != nil {
			break
		}
		articles = append(articles, a)
	}
	return
}

func FindArticlesNumByCategoryId(db *sql.DB, categoryNames []string) (articleNum int, err error) {
	query := "SELECT COUNT(id) FROM articles " +
		"WHERE id IN (" +
		"SELECT article_id FROM article_category " +
		"WHERE category_id IN (" +
		"SELECT id FROM categories "

	args, whereQuery := GenArgsFromStrSlice(categoryNames)
	query += whereQuery
	query += ")) "

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

	for rows.Next() {
		if err := rows.Scan(&articleNum); err != nil {
			ScanError.SetErr(err).AppendTo(Errors)
			break
		}
	}
	return
}

// Generate where statement of query from string slice.
func GenArgsFromStrSlice(sl []string) (args []interface{}, whereQuery string) {
	const initQuery = "WHERE "
	whereQuery = initQuery
	for i, cn := range sl {
		if i != 0 {
			whereQuery += "AND "
		}
		whereQuery += "name=? "
		args = append(args, cn)
	}
	if whereQuery == initQuery {
		whereQuery = ""
	}
	return
}
