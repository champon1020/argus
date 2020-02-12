package repository

import (
	"database/sql"
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

func (article *Article) InsertArticle(tx *sql.Tx) {

}

func (article *Article) UpdateArticle(tx *sql.Tx, argsFlg uint32) {

}

func (article *Article) DeleteArticle(tx *sql.Tx, argsFlg uint32) {

}

func (article *Article) FindArticle(tx *sql.Tx, argsFlg uint32) (articles []Article) {
	args := GenArgsSlice(argsFlg, article)
	whereQuery := GenArgsQuery(argsFlg, article)
	query := "SELECT * FROM WORDS " + whereQuery + "ORDER BY ID LIMIT 10"

	rows, err := tx.Query(query, args...)
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
		articles = append(articles, a)
	}
	return
}
