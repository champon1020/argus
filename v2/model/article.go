package model

import (
	"errors"
	"time"

	"github.com/champon1020/argus/v2"
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
	ID string `json:"id"`

	// id for sorting articles
	SortedID int `json:"sortedId"`

	// article title
	Title string `json:"title"`

	// categories of article
	Categories []Category `json:"categories"`

	// date article is posted on
	CreateDate time.Time `json:"createDate"`

	// date article is updated
	UpdateDate time.Time `json:"updateDate"`

	// content of article
	Content string `json:"content"`

	// image file name
	ImageHash string `json:"imageHash"`

	// article is private or not
	Private bool `json:"isPrivate"`
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
	*a = _a[0]

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

	return nil
}

// FindPublicArticlesByCategory searches for public articles
// which belongs to the specified category id.
func (db *Database) FindPublicArticlesByCategory(a *[]Article, categoryID int, op *QueryOptions) error {
	if db.DB == nil {
		return argus.NewError(errArticleDbNil, nil)
	}

	idCtx := db.DB.Select(nil, "categoreis", "article_id").
		Where("category_id = ?", categoryID)

	ctx := db.DB.Select(a, "articles").
		Where("private = ?", false).
		WhereCtx("id IN", idCtx)

	op.apply(ctx)

	if err := ctx.Do(); err != nil {
		return argus.NewError(errArticleQueryFailed, err).
			AppendValue("query", ctx.ToSQLString())
	}

	return nil
}

// InsertArticle inserts new article.
func (db *Database) InsertArticle(a *Article) error {
	if db.TX == nil {
		return argus.NewError(errArticleTxNil, nil)
	}

	ctx := db.TX.InsertWithModel(a, "articles")

	if err := ctx.Do(); err != nil {
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

	return nil
}
