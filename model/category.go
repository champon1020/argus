package model

import (
	"errors"

	"github.com/champon1020/argus"
)

var (
	errCategoryDbNil       = errors.New("model.category: model.Database.DB is nil")
	errCategoryQueryFailed = errors.New("model.category: Failed to execute query")
)

// Category is the struct including category information.
type Category struct {
	// unique id (primary key)
	ID string `mgorm:"id" json:"id"`

	// category name
	Name string `mgorm:"name" json:"name"`
}

// FindCategories searches for categories which is included by public articles.
func (db *Database) FindCategories(c *[]Category, op *QueryOptions) error {
	if db.DB == nil {
		return argus.NewError(errCategoryDbNil, nil)
	}

	aCtx := db.DB.Select(nil, "articles", "id").
		Where("private = ?", false)

	idCtx := db.DB.Select(nil, "article_category", "category_id").
		WhereCtx("article_id IN", aCtx)

	ctx := db.DB.Select(c, "categories", "DISTINCT *").
		WhereCtx("id IN", idCtx)

	op.apply(ctx)

	if err := ctx.Do(); err != nil {
		return argus.NewError(errCategoryQueryFailed, err).
			AppendValue("query", ctx.ToSQLString())
	}

	return nil
}

// FindCategoriesByArticleID searches for article categories by article id.
func (db *Database) FindCategoriesByArticleID(c *[]Category, articleID string) error {
	if db.DB == nil {
		return argus.NewError(errCategoryDbNil, nil)
	}

	aCtx := db.DB.Select(nil, "article_category", "category_id").
		Where("article_id = ?", articleID)

	ctx := db.DB.Select(c, "categories").
		WhereCtx("id IN", aCtx)

	if err := ctx.Do(); err != nil {
		return argus.NewError(errCategoryQueryFailed, err).
			AppendValue("query", ctx.ToSQLString())
	}

	return nil
}
