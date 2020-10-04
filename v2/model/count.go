package model

import (
	"errors"

	"github.com/champon1020/argus/v2"
)

var (
	errCountDbNil       = errors.New("model.count: model.Database.DB is nil")
	errCountQueryFailed = errors.New("model.count: Failed to execute query")
	errCountNoResult    = errors.New("model.count: Query result is nothing")
)

// Count containes the number of columns.
type Count struct {
	Value int `mgorm:"count"`
}

// CountArticles counts the number of articles.
func (db *Database) CountArticles(cnt *int) error {
	if db.DB == nil {
		return argus.NewError(errCountDbNil, nil)
	}

	var c []Count
	ctx := db.DB.Select(&c, "articles", "COUNT(*) AS count")

	if err := ctx.Do(); err != nil {
		return argus.NewError(errCountQueryFailed, err).
			AppendValue("query", ctx.ToSQLString())
	}

	if len(c) == 0 {
		return argus.NewError(errCountNoResult, nil)
	}
	*cnt = c[0].Value

	return nil
}

// CountDrafts counts the number of drafts.
func (db *Database) CountDrafts(cnt *Count) error {
	if db.DB == nil {
		return argus.NewError(errCountDbNil, nil)
	}

	var c []Count
	ctx := db.DB.Select(&c, "drafts", "COUNT(*) AS count")

	if err := ctx.Do(); err != nil {
		return argus.NewError(errCountQueryFailed, err).
			AppendValue("query", ctx.ToSQLString())
	}

	if len(c) == 0 {
		return argus.NewError(errCountNoResult, nil)
	}
	*cnt = c[0]

	return nil
}
