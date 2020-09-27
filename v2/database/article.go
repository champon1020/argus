package database

import (
	"time"
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
func (db *Database) FindArticleByID(a *[]Article, id string) error {
	ctx := db.DB.Select(a, "articles").
		Where("id = ?", id)

	return ctx.Do()
}

// FindPublicArticlesGeSortedID searches for public articles
// whose sorted id is greater than and equal to the specified
// sortedID integer.
func (db *Database) FindPublicArticlesGeSortedID(a *[]Article, sortedID int, op *QueryOptions) error {
	ctx := db.DB.Select(a, "articles").
		Where("private = ?", false).
		Where("sorted_id >= ?", sortedID)

	op.apply(ctx)

	return ctx.Do()
}

// FindAllArticles searches for all articles.
func (db *Database) FindAllArticles(a *[]Article, op *QueryOptions) error {
	ctx := db.DB.Select(a, "articles")

	op.apply(ctx)

	return ctx.Do()
}

// FindPublicArticles searches for public articles.
func (db *Database) FindPublicArticles(a *[]Article, op *QueryOptions) error {
	ctx := db.DB.Select(a, "articles").
		Where("private = ?", false)

	op.apply(ctx)

	return ctx.Do()
}

// FindPublicArticlesByTitle searches for public articles
// whose title is part of the specified title string.
func (db *Database) FindPublicArticlesByTitle(a *[]Article, title string, op *QueryOptions) error {
	ctx := db.DB.Select(a, "articles").
		Where("private = ?", false).
		Where("title LIKE %?%", title)

	op.apply(ctx)

	return ctx.Do()
}

// FindPublicArticlesByCategory searches for public articles
// which belongs to the specified category id.
func (db *Database) FindPublicArticlesByCategory(a *[]Article, categoryID int, op *QueryOptions) error {
	idCtx := db.DB.Select(nil, "categoreis", "article_id").
		Where("category_id = ?", categoryID)

	ctx := db.DB.Select(a, "articles").
		Where("private = ?", false).
		WhereCtx("id IN", idCtx)

	op.apply(ctx)

	return ctx.Do()
}

// InsertArticle inserts new article.
func (db *Database) InsertArticle(a *Article) error {
	ctx := db.TX.InsertWithModel(a, "articles")
	return ctx.Do()
}

// UpdateArticle updates the article contents.
func (db *Database) UpdateArticle(a *Article) error {
	ctx := db.TX.UpdateWithModel(a, "articles")
	return ctx.Do()
}
