package model

import (
	"errors"

	"github.com/champon1020/argus"
	"github.com/champon1020/minigorm"
)

var (
	errArticleCategoryQueryFailed = errors.New("model.article_category: Failed to execute query")
)

// insertArticleCategory inserts new pair of article and category id.
func insertArticleCategory(tx *minigorm.TX, articleID string, categoryID string) error {
	ctx := tx.Insert("article_category").
		AddColumn("article_id", articleID).
		AddColumn("category_id", categoryID)

	if err := ctx.DoTx(); err != nil {
		return argus.NewError(errArticleCategoryQueryFailed, err).
			AppendValue("query", ctx.ToSQLString())
	}

	return nil
}

// deleteArticleCategoryByArticleID deletes the pair of article and category id
// by article id.
func deleteArticleCategoryByArticleID(tx *minigorm.TX, articleID string) error {
	ctx := tx.Delete("article_category").
		Where("article_id = ?", articleID)

	if err := ctx.DoTx(); err != nil {
		return argus.NewError(errArticleCategoryQueryFailed, err).
			AppendValue("query", ctx.ToSQLString())
	}

	return nil
}
