package model

import (
	"errors"
	"sync"
	"time"

	"github.com/champon1020/argus"
	mgorm "github.com/champon1020/minigorm"
)

var (
	errArticleDbNil       = errors.New("model.article: model.Database.DB is nil")
	errArticleTxNil       = errors.New("model.article: model.Database.TX is nil")
	errArticleQueryFailed = errors.New("model.article: Failed to execute query")
	errArticleNoResult    = errors.New("model.article: Query result is nothing")
)

// Article is the struct including article information.
type Article struct {
	// unique id (primary key)
	ID string `mgorm:"id" json:"id"`

	// id for sorting articles
	SortedID int `mgorm:"sorted_id" json:"sortedId"`

	// article title
	Title string `mgorm:"title" json:"title"`

	// categories of article
	Categories []Category `mgorm:"categories" json:"categories"`

	// date article is posted on
	CreateDate time.Time `mgorm:"create_date" json:"createDate"`

	// date article is updated
	UpdateDate time.Time `mgorm:"update_date" json:"updateDate"`

	// content of article
	Content string `mgorm:"content" json:"content"`

	// image file name
	ImageHash string `mgorm:"image_hash" json:"imageHash"`

	// article is private or not
	Private bool `mgorm:"is_private" json:"isPrivate"`
}

func (db *Database) setCategoriesToArticle(a *[]Article) error {
	var err error

	wg := new(sync.WaitGroup)
	for i := 0; i < len(*a); i++ {
		wg.Add(1)
		go func(v *Article) {
			defer wg.Done()
			if e := db.FindCategoriesByArticleID(&v.Categories, v.ID); e != nil {
				err = e
			}
		}(&(*a)[i])
	}

	wg.Wait()
	if err != nil {
		return err
	}

	return nil
}

// FindArticleByID searched for the article
// whose id is the specified id string.
func (db *Database) FindArticleByID(a *Article, id string) error {
	if db.DB == nil {
		return argus.NewError(errArticleDbNil, nil)
	}

	var _a []Article
	ctx := db.DB.Select(&_a, "articles").
		Where("id = ?", id)

	if err := ctx.Do(); err != nil {
		return argus.NewError(errArticleQueryFailed, err).
			AppendValue("query", ctx.ToSQLString())
	}

	if len(_a) == 0 {
		return argus.NewError(errArticleNoResult, nil)
	}

	// Get categories by article id.
	if err := db.setCategoriesToArticle(&_a); err != nil {
		return err
	}

	*a = _a[0]

	return nil
}

// FindAllArticles searches for all articles.
func (db *Database) FindAllArticles(a *[]Article, op *QueryOptions) error {
	if db.DB == nil {
		return argus.NewError(errArticleDbNil, nil)
	}

	ctx := db.DB.Select(a, "articles")
	op.apply(ctx)

	if err := ctx.Do(); err != nil {
		return argus.NewError(errArticleQueryFailed, err).
			AppendValue("query", ctx.ToSQLString())
	}

	// Get categories by article id.
	if err := db.setCategoriesToArticle(a); err != nil {
		return err
	}

	return nil
}

// FindPublicArticles searches for public articles.
func (db *Database) FindPublicArticles(a *[]Article, op *QueryOptions) error {
	if db.DB == nil {
		return argus.NewError(errArticleDbNil, nil)
	}

	ctx := db.DB.Select(a, "articles").
		Where("private = ?", false)

	op.apply(ctx)

	if err := ctx.Do(); err != nil {
		return argus.NewError(errArticleQueryFailed, err).
			AppendValue("query", ctx.ToSQLString())
	}

	// Get categories by article id.
	if err := db.setCategoriesToArticle(a); err != nil {
		return err
	}

	return nil
}

// FindPublicArticlesGeSortedID searches for public articles
// whose sorted id is greater than and equal to the specified
// sortedID integer.
func (db *Database) FindPublicArticlesGeSortedID(a *[]Article, sortedID int, op *QueryOptions) error {
	if db.DB == nil {
		return argus.NewError(errArticleDbNil, nil)
	}

	ctx := db.DB.Select(a, "articles").
		Where("private = ?", false).
		Where("sorted_id >= ?", sortedID)

	op.apply(ctx)

	if err := ctx.Do(); err != nil {
		return argus.NewError(errArticleQueryFailed, err).
			AppendValue("query", ctx.ToSQLString())
	}

	// Get categories by article id.
	if err := db.setCategoriesToArticle(a); err != nil {
		return err
	}

	return nil
}

// FindPublicArticlesByTitle searches for public articles
// whose title is part of the specified title string.
func (db *Database) FindPublicArticlesByTitle(a *[]Article, title string, op *QueryOptions) error {
	ctx := db.DB.Select(a, "articles").
		Where("private = ?", false).
		Where("title LIKE %?%", title)

	op.apply(ctx)

	if err := ctx.Do(); err != nil {
		return argus.NewError(errArticleQueryFailed, err).
			AppendValue("query", ctx.ToSQLString())
	}

	// Get categories by article id.
	if err := db.setCategoriesToArticle(a); err != nil {
		return err
	}

	return nil
}

// FindPublicArticlesByCategory searches for public articles
// which belongs to the specified category id.
func (db *Database) FindPublicArticlesByCategory(a *[]Article, categoryID string, op *QueryOptions) error {
	if db.DB == nil {
		return argus.NewError(errArticleDbNil, nil)
	}

	idCtx := db.DB.Select(nil, "article_category", "article_id").
		Where("category_id = ?", categoryID)

	ctx := db.DB.Select(a, "articles").
		Where("private = ?", false).
		WhereCtx("id IN", idCtx)

	op.apply(ctx)

	if err := ctx.Do(); err != nil {
		return argus.NewError(errArticleQueryFailed, err).
			AppendValue("query", ctx.ToSQLString())
	}

	// Get categories by article id.
	if err := db.setCategoriesToArticle(a); err != nil {
		return err
	}

	return nil
}

// RegisterArticle registers new article.
// This function inserts new article and new categories to
// each tables on transaction process.
func (db *Database) RegisterArticle(a *Article) error {
	err := db.DB.Transact(func(tx mgorm.TX) error {
		doneArt := make(chan bool)
		flgCate := true
		errc := make(chan error)

		// Insert article to database.
		go func(v *Article) {
			if err := insertArticle(&tx, v); err != nil {
				doneArt <- false
				errc <- err
				return
			}
			doneArt <- true
		}(a)

		// Insert category and pair of article and category id to database.
		wg := new(sync.WaitGroup)
		for i := 0; i < len(a.Categories); i++ {
			wg.Add(1)
			go func(c *Category) {
				defer wg.Done()
				if err := insertCategories(&tx, c); err != nil {
					flgCate = false
					errc <- err
					return
				}

				if err := insertArticleCategory(&tx, a.ID, c.ID); err != nil {
					flgCate = false
					errc <- err
					return
				}
			}(&a.Categories[i])
		}

		wg.Wait()

		if !<-doneArt || !flgCate {
			return <-errc
		}

		return nil
	})

	return err
}

// insertArticle inserts new article.
func insertArticle(tx *mgorm.TX, a *Article) error {
	ctx := tx.Insert("articles").
		AddColumn("id", getNewID()).
		AddColumn("title", a.Title).
		AddColumn("create_date", time.Now()).
		AddColumn("update_date", time.Now()).
		AddColumn("content", a.Content).
		AddColumn("image_hash", a.ImageHash).
		AddColumn("private", a.Private)

	if err := ctx.DoTx(); err != nil {
		return argus.NewError(errArticleQueryFailed, err).
			AppendValue("query", ctx.ToSQLString())
	}

	return nil
}

// UpdateArticle updates the article contents.
func (db *Database) UpdateArticle(a *Article) error {
	if db.TX == nil {
		return argus.NewError(errArticleTxNil, nil)
	}

	ctx := db.TX.UpdateWithModel(a, "articles")

	if err := ctx.Do(); err != nil {
		return argus.NewError(errArticleQueryFailed, err).
			AppendValue("query", ctx.ToSQLString())
	}

	// TODO: implement process to update category.

	return nil
}
