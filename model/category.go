package model

import (
	"errors"

	"github.com/champon1020/argus"
	"github.com/champon1020/minigorm"
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

// FindPublicCategories searches for categories which is included by public articles.
func (db *Database) FindPublicCategories(c *[]Category, op *QueryOptions) error {
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

// insertCategory inserts new category.
func insertCategories(tx *minigorm.TX, c *Category) error {
	if err := assignCategoryIDIfExist(tx, c); err != nil {
		return err
	} else if c.ID != "" {
		// Category has already existed.
		return nil
	}

	// Generate new category id.
	c.ID = getNewID()

	ctx := tx.InsertWithModel(c, "categories")
	if err := ctx.DoTx(); err != nil {
		return argus.NewError(errCategoryQueryFailed, err).
			AppendValue("query", ctx.ToSQLString())
	}

	return nil
}

// assignCategoryIDIfExist assigns id to category object.
// If the category is exist in database table, get id from database.
func assignCategoryIDIfExist(tx *minigorm.TX, c *Category) error {
	res := []struct {
		ID string `mgorm:"id"`
	}{}

	ctx := tx.Select(&res, "categories", "id").
		Where("name = ?", c.Name)

	if err := ctx.DoTx(); err != nil {
		return argus.NewError(errCategoryQueryFailed, err).
			AppendValue("query", ctx.ToSQLString())
	}

	if len(res) > 0 {
		c.ID = res[0].ID
	}

	return nil
}
