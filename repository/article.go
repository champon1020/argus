package repository

import (
	"database/sql"
	"sync"
	"time"
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
	Id          int        `json:"id"`
	Title       string     `json:"title"`
	Categories  []Category `json:"categories"`
	CreateDate  time.Time  `json:"createDate"`
	UpdateDate  time.Time  `json:"updateDate"`
	ContentHash string     `json:"contentHash"`
	ImageHash   string     `json:"imageHash"`
	Private     bool       `json:"private"`
}

func (article *Article) InsertArticle(tx *sql.Tx) (err error) {
	cmd := "INSERT INTO articles " +
		"(id, title, create_date, update_date, content_hash, image_hash, private)" +
		"VALUES (?, ?, ?, ?, ?, ?, ?)"

	if _, err := tx.Exec(cmd,
		article.Id,
		article.Title,
		article.CreateDate,
		article.UpdateDate,
		article.ContentHash,
		article.ImageHash,
		article.Private,
	); err != nil {
		CmdError.SetErr(err).AppendTo(Errors)
	}
	return
}

// Insert column to article_category table.
func (article *Article) InsertArticleCategory(tx *sql.Tx) (err error) {
	cmd := "INSERT INTO article_category (article_id, category_id) " +
		"VALUES (?, ?)"

	wg := new(sync.WaitGroup)
	for _, c := range article.Categories {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err = tx.Exec(cmd, article.Id, c.Id); err != nil {
				CmdError.SetErr(err).AppendTo(Errors)
			}
		}()
	}
	wg.Wait()
	return
}

func (article *Article) UpdateArticle(tx *sql.Tx) (err error) {
	cmd := "UPDATE articles " +
		"SET title=?, update_date=?, content_hash=?, image_hash=?, private=? " +
		"WHERE id=?"

	if _, err = tx.Exec(cmd,
		article.Title,
		article.UpdateDate,
		article.ContentHash,
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

// Remove column which of article_id is equal to object from article_category table.
func (article *Article) DeleteArticleCategoryByArticle(tx *sql.Tx) (err error) {
	cmd := "DELETE FROM article_category WHERE article_id=?"
	if _, err = tx.Exec(cmd, article.Id); err != nil {
		CmdError.SetErr(err).AppendTo(Errors)
	}
	return
}

// Remove column that both of article_id and category_id is equal to object.
func (article *Article) DeleteArticleCategoryByBoth(tx *sql.Tx) (err error) {
	cmd := "DELETE FROM article_category WHERE article_id=? AND category_id=?"

	wg := new(sync.WaitGroup)
	for _, c := range article.Categories {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err = tx.Exec(cmd, article.Id, c.Id); err != nil {
				CmdError.SetErr(err).AppendTo(Errors)
			}
		}()
	}
	wg.Wait()
	return
}

// ArgFlg determines where statement's arguments.
// For Example, 'argsFlg = 0101' means
// it includes first and third fields of objects in where statement.
func (article *Article) FindArticle(db *sql.DB, argsFlg uint32, offset int) (articles []Article, err error) {
	args := GenArgsSlice(argsFlg, article, offset)
	whereQuery, limitQuery := GenArgsQuery(argsFlg, article)
	query := "SELECT * FROM articles " + whereQuery +
		"ORDER BY id DESC " + limitQuery +
		"OFFSET ?"

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

func (article *Article) FindArticlesNum(db *sql.DB, argsFlg uint32) (articleNum int, err error) {
	args := GenArgsSlice(argsFlg, article, -1)
	whereQuery, _ := GenArgsQuery(argsFlg, article)
	query := "SELECT COUNT(id) FROM articles " + whereQuery

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
