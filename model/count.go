package model

import (
	"errors"

	"github.com/champon1020/argus"
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

// CountAllArticles counts the number of all articles.
func (db *Database) CountAllArticles(cnt *int, op *QueryOptions) error {
	if db.DB == nil {
		return argus.NewError(errCountDbNil, nil)
	}

	var c []Count
	ctx := db.DB.Select(&c, "articles", "COUNT(*) AS count")
	op.apply(ctx)

	if err := ctx.Do(); err != nil {
		return argus.NewError(errCountQueryFailed, err).
			AppendValue("query", ctx.ToSQLString())
	}

	if err := assignCount(cnt, &c); err != nil {
		return err
	}

	return nil
}

// CountPublicArticles counts the number of public articles.
func (db *Database) CountPublicArticles(cnt *int, op *QueryOptions) error {
	if db.DB == nil {
		return argus.NewError(errCountDbNil, nil)
	}

	var c []Count
	ctx := db.DB.Select(&c, "articles", "COUNT(*) AS count").
		Where("private = ?", false)
	op.apply(ctx)

	if err := ctx.Do(); err != nil {
		return argus.NewError(errCountQueryFailed, err).
			AppendValue("query", ctx.ToSQLString())
	}

	if err := assignCount(cnt, &c); err != nil {
		return err
	}

	return nil
}

// CountPublicArticlesByTitle counts the number of articles with title.
func (db *Database) CountPublicArticlesByTitle(cnt *int, title string, op *QueryOptions) error {
	if db.DB == nil {
		return argus.NewError(errCountDbNil, nil)
	}

	var c []Count
	ctx := db.DB.Select(&c, "articles", "COUNT(*) AS count").
		Where("private = ?", false).
		Where("title LIKE %?%", title)

	op.apply(ctx)

	if err := ctx.Do(); err != nil {
		return argus.NewError(errCountQueryFailed, err).
			AppendValue("query", ctx.ToSQLString())
	}

	if err := assignCount(cnt, &c); err != nil {
		return err
	}

	return nil
}

// CountPublicArticlesByCategory counts the number of articles with category.
func (db *Database) CountPublicArticlesByCategory(cnt *int, categoryID string, op *QueryOptions) error {
	if db.DB == nil {
		return argus.NewError(errCountDbNil, nil)
	}

	idCtx := db.DB.Select(nil, "article_category", "article_id").
		Where("category_id = ?", categoryID)

	var c []Count
	ctx := db.DB.Select(&c, "articles", "COUNT(*) AS count").
		Where("private = ?", false).
		WhereCtx("id IN", idCtx)

	op.apply(ctx)

	if err := ctx.Do(); err != nil {
		return argus.NewError(errArticleQueryFailed, err).
			AppendValue("query", ctx.ToSQLString())
	}

	if err := assignCount(cnt, &c); err != nil {
		return err
	}

	return nil
}

// CountDrafts counts the number of drafts.
func (db *Database) CountDrafts(cnt *int, op *QueryOptions) error {
	if db.DB == nil {
		return argus.NewError(errCountDbNil, nil)
	}

	var c []Count
	ctx := db.DB.Select(&c, "drafts", "COUNT(*) AS count")
	op.apply(ctx)

	if err := ctx.Do(); err != nil {
		return argus.NewError(errCountQueryFailed, err).
			AppendValue("query", ctx.ToSQLString())
	}

	if err := assignCount(cnt, &c); err != nil {
		return err
	}

	return nil
}

func assignCount(cnt *int, c *[]Count) error {
	if len(*c) == 0 {
		return argus.NewError(errCountNoResult, nil)
	}

	*cnt = (*c)[0].Value
	return nil
}
