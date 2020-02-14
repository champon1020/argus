package repository

import (
	"database/sql"
	"sync"
	"time"
)

type Article struct {
	Id         int
	Title      string
	Categories []Category
	CreateDate time.Time
	UpdateDate time.Time
	ContentUrl string
	ImageUrl   string
	Private    bool
}

func (article *Article) InsertArticle(tx *sql.Tx) (err error) {
	cmd := "INSERT INTO articles " +
		"(id, title, create_date, update_date, content_url, image_url, private)" +
		"VALUES (?, ?, ?, ?, ?, ?, ?)"

	_, err = tx.Exec(cmd,
		article.Id,
		article.Title,
		article.CreateDate,
		article.UpdateDate,
		article.ContentUrl,
		article.ImageUrl,
		article.Private)

	if err != nil {
		logger.ErrorPrintf(err)
	}
	return
}

func (article *Article) InsertArticleCategory(tx *sql.Tx) (err error) {
	cmd := "INSERT INTO article_category (article_id, category_id) " +
		"VALUES (?, ?)"

	wg := new(sync.WaitGroup)
	for _, c := range article.Categories {
		wg.Add(1)
		go func(c Category) {
			defer wg.Done()

			_, err = tx.Exec(cmd, article.Id, c.Id)
			if err != nil {
				logger.ErrorPrintf(err)
			}
		}(c)
	}
	wg.Wait()
	return
}

func (article *Article) UpdateArticle(tx *sql.Tx) (err error) {
	cmd := "UPDATE articles " +
		"SET title=?, create_date=?, update_date=?, content_url=?, image_url=?, private=?" +
		"WHERE id=?"

	_, err = tx.Exec(cmd,
		article.Title,
		article.CreateDate,
		article.UpdateDate,
		article.ContentUrl,
		article.ImageUrl,
		article.Private,
		article.Id)

	if err != nil {
		logger.ErrorPrintf(err)
	}
	return
}

func (article *Article) DeleteArticle(tx *sql.Tx) (err error) {
	cmd := "DELETE FROM articles WHERE id=?"
	_, err = tx.Exec(cmd, article.Id)
	if err != nil {
		logger.ErrorPrintf(err)
	}
	return
}

func (article *Article) DeleteArticleCategoryByArticle(tx *sql.Tx) (err error) {
	cmd := "DELETE FROM article_category WHERE article_id=?"
	_, err = tx.Exec(cmd, article.Id)
	if err != nil {
		logger.ErrorPrintf(err)
	}
	return
}

func (article *Article) DeleteArticleCategoryByBoth(tx *sql.Tx) (err error) {
	cmd := "DELETE FROM article_category WHERE article_id=? AND category_id=?"

	wg := new(sync.WaitGroup)
	for _, c := range article.Categories {
		wg.Add(1)
		go func(c Category) {
			defer wg.Done()

			_, err = tx.Exec(cmd, article.Id, c.Id)
			if err != nil {
				logger.ErrorPrintf(err)
			}
		}(c)
	}
	wg.Wait()
	return
}

func (article *Article) FindArticle(db *sql.DB, argsFlg uint32) (articles []Article) {
	args := GenArgsSlice(argsFlg, article)
	whereQuery := GenArgsQuery(argsFlg, article)
	query := "SELECT * FROM articles " + whereQuery + "ORDER BY id LIMIT 10"

	rows, err := db.Query(query, args...)
	defer func() {
		if err := rows.Close(); err != nil {
			logger.ErrorPrintf(err)
		}
	}()

	if err != nil {
		logger.ErrorPrintf(err)
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
		a.Categories = a.FindArticleCategory(db)
		articles = append(articles, a)
	}
	return
}

func (article *Article) FindArticleCategory(db *sql.DB) (categories []Category) {
	query := "SELECT * FROM categories " +
		"WHERE id IN (" +
		"SELECT category_id FROM article_category " +
		"WHERE article_id=?)"

	rows, err := db.Query(query, article.Id)
	defer func() {
		if err := rows.Close(); err != nil {
			logger.ErrorPrintf(err)
		}
	}()

	if err != nil {
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
