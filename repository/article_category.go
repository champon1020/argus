package repository

import "database/sql"

func FindArticleByCategoryId(db *sql.DB, categoryNames []string, argsFlg uint32) (articles []Article, err error) {
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
	query += ")) "

	args = append(args, GenArgsSlice(argsFlg, Article{})...)
	_, limitQuery := GenArgsQuery(argsFlg, Article{})
	query += limitQuery

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
