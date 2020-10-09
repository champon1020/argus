package repo

import (
	"database/sql"

	"github.com/champon1020/argus/v1/service"
)

type ArticleCategory struct {
	ArticleId  string
	CategoryId string
}

// Insert column to article_category table.
func (ac *ArticleCategory) InsertArticleCategory(tx *sql.Tx) (err error) {
	cmd := "INSERT INTO article_category (article_id, category_id) " +
		"SELECT ?,? WHERE NOT EXISTS (" +
		"SELECT * FROM article_category WHERE article_id=? AND category_id=?)"

	args := []interface{}{
		ac.ArticleId, ac.CategoryId,
		ac.ArticleId, ac.CategoryId,
	}
	if _, err = tx.Exec(cmd, args...); err != nil {
		CmdError.
			SetErr(err).
			SetValues("cmd", cmd).
			SetValues("args", args).
			AppendTo(Errors)
		return
	}
	return
}

func DeleteArticleCategory(tx *sql.Tx, option *service.QueryOption) (err error) {
	args := service.GenArgsSlice(*option)
	query := service.GenArgsQuery(*option)
	cmd := "DELETE FROM article_category " + query
	if _, err = tx.Exec(cmd, args...); err != nil {
		CmdError.
			SetErr(err).
			SetValues("cmd", cmd).
			SetValues("args", args).
			AppendTo(Errors)
	}
	return
}

func FindArticleByCategoryId(db *sql.DB, option *service.QueryOption) (articles []Article, err error) {
	query := "SELECT * FROM articles " +
		"WHERE id IN (" +
		"SELECT article_id FROM article_category " +
		"WHERE category_id IN (" +
		"SELECT id FROM categories "

	query += service.GenWhereQuery(*option) + ")) "
	query += service.GenOrderQuery(*option)
	query += service.GenLimitQuery(*option)

	args := service.GenArgsSlice(*option)

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
			&a.SortedId,
			&a.Title,
			&a.CreateDate,
			&a.UpdateDate,
			&a.Content,
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

func FindArticlesNumByCategoryId(db *sql.DB, option *service.QueryOption) (articleNum int, err error) {
	query := "SELECT COUNT(id) FROM articles " +
		"WHERE id IN (" +
		"SELECT article_id FROM article_category " +
		"WHERE category_id IN (" +
		"SELECT id FROM categories "

	query += service.GenWhereQuery(*option) + ")) "
	query += service.GenOrderQuery(*option)
	query += service.GenLimitQuery(*option)

	args := service.GenArgsSlice(*option)

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
