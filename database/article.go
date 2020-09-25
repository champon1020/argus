package database

import (
	"time"

	mgorm "github.com/champon1020/minigorm"
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
func FindArticleByID(db *mgorm.DB, a *[]Article, id string) error {
	ctx := db.Select(a, "articles").
		Where("id = ?", id)

	return ctx.Do()
}

// FindPublicArticlesGeSortedID searches for public articles
// whose sorted id is greater than and equal to the specified
// sortedID integer.
func FindPublicArticlesGeSortedID(db *mgorm.DB, a *[]Article, sortedID int, op *QueryOptions) error {
	ctx := db.Select(a, "articles").
		Where("private = ?", false).
		Where("sorted_id >= ?", sortedID)

	op.apply(ctx)

	return ctx.Do()
}

// FindAllArticles searches for all articles.
func FindAllArticles(db *mgorm.DB, a *[]Article, op *QueryOptions) error {
	ctx := db.Select(a, "articles")

	op.apply(ctx)

	return ctx.Do()
}

// FindPublicArticles searches for public articles.
func FindPublicArticles(db *mgorm.DB, a *[]Article, op *QueryOptions) error {
	ctx := db.Select(a, "articles").
		Where("private = ?", false)

	op.apply(ctx)

	return ctx.Do()
}

// FindPublicArticlesByTitle searches for public articles
// whose title is part of the specified title string.
func FindPublicArticlesByTitle(db *mgorm.DB, a *[]Article, title string, op *QueryOptions) error {
	ctx := db.Select(a, "articles").
		Where("private = ?", false).
		Where("title LIKE %?%", title)

	op.apply(ctx)

	return ctx.Do()
}

// FindPublicArticlesByCategory searches for public articles
// which belongs to the specified category.
func FindPublicArticlesByCategory(db *mgorm.DB, a *[]Article, category Category, op *QueryOptions) error {
	idCtx := db.Select(nil, "categoreis", "article_id").
		Where("category_id = ?", category.ID)

	ctx := db.Select(a, "articles").
		Where("private = ?", false).
		WhereCtx("id IN", idCtx)

	op.apply(ctx)

	return ctx.Do()
}

// FindNewPublicArticles searches for new public articles.
func FindNewPublicArticles(db *mgorm.DB, a *[]Article) error {
	op := &QueryOptions{
		Limit:   5,
		OrderBy: "sorted_id",
		Desc:    true,
	}
	return FindPublicArticles(db, a, op)
}
