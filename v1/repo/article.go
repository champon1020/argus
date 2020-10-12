package repo

import (
	"database/sql"
	"time"

	"github.com/champon1020/argus/v1/service"
)

// Id: primary key
// Title: article title
// Categories: categories of article
// CreateDate: created date
// UpdateDate: last updated date
// ContentHash: content file name (html file)
// ImageHash: image file name
// private: this article is whether public or not
type Article struct {
	Id         string     `json:"id"`
	SortedId   int        `json:"sortedId"`
	Title      string     `json:"title"`
	Categories []Category `json:"categories"`
	CreateDate time.Time  `json:"createDate"`
	UpdateDate time.Time  `json:"updateDate"`
	Content    string     `json:"content"`
	ImageHash  string     `json:"imageHash"`
	Private    bool       `json:"isPrivate"`
}

func (article *Article) InsertArticle(tx *sql.Tx) (err error) {
	cmd := "INSERT INTO articles " +
		"(id, title, create_date, update_date, content, image_hash, private)" +
		"VALUES (?, ?, ?, ?, ?, ?, ?)"

	if _, err := tx.Exec(cmd,
		article.Id,
		article.Title,
		article.CreateDate,
		article.UpdateDate,
		article.Content,
		article.ImageHash,
		article.Private,
	); err != nil {
		CmdError.SetErr(err).AppendTo(Errors)
	}
	return
}

func (article *Article) UpdateArticle(tx *sql.Tx) (err error) {
	cmd := "UPDATE articles " +
		"SET title=?, update_date=?, content=?, image_hash=?, private=? " +
		"WHERE id=?"

	if _, err = tx.Exec(cmd,
		article.Title,
		article.UpdateDate,
		article.Content,
		article.ImageHash,
		article.Private,
		article.Id,
	); err != nil {
		CmdError.SetErr(err).AppendTo(Errors)
	}
	return
}

func (article *Article) DeleteArticle(tx *sql.Tx) (err error) {
	cmd := "DELETE FROM articles WHERE id=?"
	if _, err = tx.Exec(cmd, article.Id); err != nil {
		CmdError.SetErr(err).AppendTo(Errors)
	}
	return
}

// ArgFlg determines where statement's arguments.
// For Example, 'argsMask = 0101' means
// it includes first and third fields of objects in where statement.
func FindArticle(db *sql.DB, option *service.QueryOption) (articles []Article, err error) {
	args := service.GenArgsSlice(*option)
	argsQuery := service.GenArgsQuery(*option)
	query := "SELECT * FROM articles " + argsQuery

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

// Find categories from categories table which of column of article_id is equal to object.
func (article *Article) FindCategoryByArticleId(db *sql.DB) (categories []Category, err error) {
	query := "SELECT * FROM categories " +
		"WHERE id IN (" +
		"SELECT category_id FROM article_category " +
		"WHERE article_id=?) ORDER BY name"

	var rows *sql.Rows
	defer RowsClose(rows)
	if rows, err = db.Query(query, article.Id); err != nil || rows == nil {
		QueryError.
			SetErr(err).
			SetValues("query", query).
			SetValues("args", article.Id).
			AppendTo(Errors)
		return
	}

	var c Category
	for rows.Next() {
		if err := rows.Scan(
			&c.Id,
			&c.Name); err != nil {
			ScanError.SetErr(err).AppendTo(Errors)
		}
		categories = append(categories, c)
	}
	return
}

func FindArticlesNum(db *sql.DB, option *service.QueryOption) (articleNum int, err error) {
	args := service.GenArgsSlice(*option)
	argsQuery := service.GenArgsQuery(*option)
	query := "SELECT COUNT(id) FROM articles " + argsQuery

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
		}
	}
	return
}
