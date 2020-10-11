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
	cmd := "INSERT INTO article_category (article_id, category_id) " +
		"SELECT ?,? WHERE NOT EXISTS " +
		"(SELECT * FROM article_category WHERE article_id=? AND category_id=?)"

	if err := tx.RawExec(cmd, articleID, categoryID, articleID, categoryID); err != nil {
		return argus.NewError(errArticleCategoryQueryFailed, err).
			AppendValue("query", cmd)
	}

	return nil
}
