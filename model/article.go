package model

import (
	"errors"
	"sync"
	"time"

	"github.com/champon1020/argus"
	"github.com/champon1020/minigorm"
)

var (
	errArticleDbNil       = errors.New("model.article: model.Database.DB is nil")
	errArticleTxNil       = errors.New("model.article: model.Database.TX is nil")
	errArticleQueryFailed = errors.New("model.article: Failed to execute query")
	errArticleNoResult    = errors.New("model.article: Query result is nothing")
	errFailedBeginTx      = errors.New("model.article: Failed to begin transaction")
)

// Article is the struct including article information.
type Article struct {
	// unique id (primary key)
	ID string `mgorm:"id" json:"id"`

	// article title
	Title string `mgorm:"title" json:"title"`

	// categories of article
	Categories []Category `mgorm:"categories" json:"categories"`

	// date article is posted on
	CreatedDate time.Time `mgorm:"created_date" json:"createdDate"`

	// date article is updated
	UpdatedDate time.Time `mgorm:"updated_date" json:"updatedDate"`

	// content of article
	Content string `mgorm:"content" json:"content"`

	// image file name
	ImageHash string `mgorm:"image_name" json:"imageName"`

	// article is private or not
	Private bool `mgorm:"private" json:"private"`
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

	return err
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

// FindPublicArticleLeID searches for public article
// whose id is less than or equal to argument's id.
func (db *Database) FindPublicArticleLeID(a *[]Article, id string, op *QueryOptions) error {
	if db.DB == nil {
		return argus.NewError(errArticleDbNil, nil)
	}

	ctx := db.DB.Select(a, "articles").
		Where("private = ?", false).
		Where("id <= ?", id)

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

// FindPublicArticleGeID search for public article
// whose id is greater than and equal to argument's id.
func (db *Database) FindPublicArticleGeID(a *[]Article, id string, op *QueryOptions) error {
	if db.DB == nil {
		return argus.NewError(errArticleDbNil, nil)
	}

	ctx := db.DB.Select(a, "articles").
		Where("private = ?", false).
		Where("id >= ?", id)

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
func (db *Database) RegisterArticle(a *Article, draftID string) error {
	// Create transaction instance.
	tx, err := db.DB.NewTX()
	if err != nil {
		return argus.NewError(errFailedBeginTx, err)
	}

	// Generate new article id.
	a.ID = GetNewID(TypeArticle)

	err = tx.Transact(func(tx *minigorm.TX) error {
		flgCate := true
		errc := make(chan error, 3)

		// If this article's draft has already existed, delete it at first.
		if draftID != "" {
			deleteDraft(tx, draftID)
		}

		// Insert article to database.
		if err := insertArticle(tx, a); err != nil {
			return err
		}

		// Insert categories and pairs of article and category id to database.
		wg := new(sync.WaitGroup)
		for i := 0; i < len(a.Categories); i++ {
			wg.Add(1)
			go func(c *Category) {
				defer wg.Done()
				if err := insertCategories(tx, c); err != nil {
					flgCate = false
					errc <- err
					return
				}

				if err := insertArticleCategory(tx, a.ID, c.ID); err != nil {
					flgCate = false
					errc <- err
					return
				}
			}(&a.Categories[i])
		}

		wg.Wait()

		if !flgCate {
			return <-errc
		}

		return nil
	})

	return err
}

// insertArticle inserts new article.
func insertArticle(tx *minigorm.TX, a *Article) error {
	ctx := tx.Insert("articles").
		AddColumn("id", a.ID).
		AddColumn("title", a.Title).
		AddColumn("created_date", time.Now()).
		AddColumn("updated_date", time.Now()).
		AddColumn("content", a.Content).
		AddColumn("image_name", a.ImageHash).
		AddColumn("private", boolToInt(a.Private))

	if err := ctx.DoTx(); err != nil {
		return argus.NewError(errArticleQueryFailed, err).
			AppendValue("query", ctx.ToSQLString())
	}

	return nil
}

// UpdateArticle updates existed article.
func (db *Database) UpdateArticle(a *Article) error {
	// Create transaction instance.
	tx, err := db.DB.NewTX()
	if err != nil {
		return argus.NewError(errFailedBeginTx, err)
	}

	err = tx.Transact(func(tx *minigorm.TX) error {
		doneUpdate := make(chan bool, 1)
		doneDelete := make(chan bool, 1)
		errc := make(chan error, 3)

		// Update article on database.
		go func() {
			defer close(doneUpdate)
			if err := updateArticle(tx, a); err != nil {
				errc <- err
				doneUpdate <- false
				return
			}
			doneUpdate <- true
		}()

		// Delete pairs of article and category id.
		go func() {
			defer close(doneDelete)
			if err := deleteArticleCategoryByArticleID(tx, a.ID); err != nil {
				errc <- err
				doneDelete <- false
				return
			}
			doneDelete <- true
		}()

		if !<-doneUpdate || !<-doneDelete {
			return <-errc
		}

		flgCate := true

		// Insert categories and pairs of article and category id to database.
		wg := new(sync.WaitGroup)
		for i := 0; i < len(a.Categories); i++ {
			wg.Add(1)
			go func(c *Category) {
				defer wg.Done()
				if err := insertCategories(tx, c); err != nil {
					flgCate = false
					errc <- err
				}

				if err := insertArticleCategory(tx, a.ID, c.ID); err != nil {
					flgCate = false
					errc <- err
				}
			}(&a.Categories[i])
		}

		wg.Wait()

		if !flgCate {
			return <-errc
		}

		// Delete categories which is not used.
		if err := deleteCategoriesNotUsed(tx); err != nil {
			return err
		}

		return nil
	})

	return err
}

// updateArticle updates the article contents.
func updateArticle(tx *minigorm.TX, a *Article) error {
	ctx := tx.Update("articles").
		AddColumn("title", a.Title).
		AddColumn("updated_date", time.Now()).
		AddColumn("content", a.Content).
		AddColumn("image_name", a.ImageHash).
		AddColumn("private", boolToInt(a.Private)).
		Where("id = ?", a.ID)

	if err := ctx.DoTx(); err != nil {
		return argus.NewError(errArticleQueryFailed, err).
			AppendValue("query", ctx.ToSQLString())
	}

	return nil
}

// UpdateArticlePrivate updates private column of article.
func (db *Database) UpdateArticlePrivate(id string, isPrivate bool) error {
	tx, err := db.DB.NewTX()
	if err != nil {
		return argus.NewError(errFailedBeginTx, err)
	}

	err = tx.Transact(func(tx *minigorm.TX) error {
		ctx := tx.Update("articles").
			AddColumn("private", boolToInt(isPrivate)).
			Where("id = ?", id)

		if err := ctx.DoTx(); err != nil {
			return argus.NewError(errArticleQueryFailed, err).
				AppendValue("query", ctx.ToSQLString())
		}

		return nil
	})

	return err
}
