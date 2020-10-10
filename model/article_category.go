package model

import (
	"errors"

	"github.com/champon1020/argus"
	mgorm "github.com/champon1020/minigorm"
)

var (
	errArticleCategoryQueryFailed = errors.New("model.article_category: Failed to execute query")
)

// insertArticleCategory inserts new pair of article and category id.
func insertArticleCategory(tx *mgorm.TX, articleID string, categoryID string) error {
	cmd := "INSERT INTO article_category (article_id, category_id) " +
		"SELECT article_id, category_id FROM dual " +
		"WHERE NOT EXISTS(SELECT * FROM article_category" +
		"WHERE article_id = ? AND category_id = ?"

	if err := tx.RawExec(cmd, articleID, categoryID); err != nil {
		return argus.NewError(errArticleCategoryQueryFailed, err).
			AppendValue("query", cmd)
	}

	return nil
}
