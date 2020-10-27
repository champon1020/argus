package model

import (
	"errors"
	"sync"

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

	// number of articles which belongs to this category
	ArticleNum int `json:"articleNum"`
}

func (db *Database) setArticleNumToCategory(c *[]Category) error {
	var err error

	wg := new(sync.WaitGroup)
	for i := 0; i < len(*c); i++ {
		wg.Add(1)
		go func(v *Category) {
			defer wg.Done()
			if e := db.CountCategoriesByPublicArticles(&v.ArticleNum, v.ID); e != nil {
				err = e
			}
		}(&(*c)[i])
	}

	wg.Wait()

	return err
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

	if err := db.setArticleNumToCategory(c); err != nil {
		return err
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

	if err := db.setArticleNumToCategory(c); err != nil {
		return err
	}

	return nil
}

// insertCategory inserts new category.
// If category has already existed, it assigns the id to struct field.
func insertCategories(tx *minigorm.TX, c *Category) error {
	if err := assignCategoryIDIfExist(tx, c); err != nil {
		return err
	} else if c.ID != "" {
		// Category has already existed.
		return nil
	}

	// Generate new category id.
	c.ID = GetNewID(TypeCategory)

	ctx := tx.Insert("categories").
		AddColumn("id", c.ID).
		AddColumn("name", c.Name)

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

func deleteCategoriesNotUsed(tx *minigorm.TX) error {
	cmd := "DELETE c FROM categories AS c " +
		"WHERE NOT EXISTS " +
		"(SELECT * FROM article_category AS ac WHERE c.id = ac.category_id)"

	if err := tx.RawExec(cmd); err != nil {
		return argus.NewError(errCategoryQueryFailed, err).
			AppendValue("query", cmd)
	}

	return nil
}
