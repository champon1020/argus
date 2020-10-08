package model

import (
	"errors"

	"github.com/champon1020/argus/v2"
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

// FindCategories searches for article categories.
func (db *Database) FindCategories(c *[]Category, op *QueryOptions) error {
	if db.DB == nil {
		return argus.NewError(errCategoryDbNil, nil)
	}

	ctx := db.DB.Select(c, "categories")
	op.apply(ctx)

	if err := ctx.Do(); err != nil {
		return argus.NewError(errCategoryQueryFailed, err).
			AppendValue("query", ctx.ToSQLString())
	}

	return nil
}
