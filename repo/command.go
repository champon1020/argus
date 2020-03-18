package repo

import (
	"database/sql"
	"errors"
	"sync"

	"github.com/champon1020/argus/service"
)

type RegisterArticleCmd func(MySQL, Article) error

/*
Flow:
	- Convert articles id.
	- Insert articles.
	- Sample new categories from articles.Categories.
	- Insert new categories.
	- Insert the pair of article_id and category_ids.
*/
func RegisterArticleCommand(mysql MySQL, article Article) (err error) {
	var newCa []Category
	if newCa, _, err = ExtractCategory(mysql.DB, article); err != nil {
		return
	}
	for _, c := range newCa {
		service.GenNewId(service.IdLen, &c.Id)
	}
	article.Categories = newCa

	var d []Draft
	option := &service.QueryOption{
		Args: []interface{}{article.ContentHash},
		Aom: map[string]service.Ope{
			"ContentHash": service.Eq,
		},
	}
	if d, err = FindDrafts(mysql.DB, option); err != nil {
		return
	}

	// Start transaction
	err = mysql.Transact(func(tx *sql.Tx) (err error) {
		if len(d) > 0 {
			d[0].DeleteDraft(tx)
		}
		if err = InsertCategories(tx, newCa); err != nil {
			return
		}
		if err = article.InsertArticle(tx); err != nil {
			return
		}
		if err = article.InsertArticleCategory(tx); err != nil {
			return
		}
		return
	})
	// End transaction

	return
}

type UpdateArticleCmd func(MySQL, Article) error

/*
Flow:
	- Get categories by article_id.
	- Classify categories to new and old.
	- Update articles.
	- Insert new categories.
	- Insert the pair of article_id and new category_ids.
	- Delete the pair of article_id and old category_ids.
*/
func UpdateArticleCommand(mysql MySQL, article Article) (err error) {
	var newCa, delCa []Category
	if newCa, delCa, err = ExtractCategory(mysql.DB, article); err != nil {
		return
	}
	for _, c := range newCa {
		service.GenNewId(service.IdLen, &c.Id)
	}

	// Start transaction
	err = mysql.Transact(func(tx *sql.Tx) (err error) {
		errCnt := 0
		wg := new(sync.WaitGroup)
		wg.Add(3)

		// update articles
		go func() {
			defer wg.Done()
			if err = article.UpdateArticle(tx); err != nil {
				errCnt++
			}
		}()

		// insert new categories
		go func() {
			defer wg.Done()
			a := Article{Id: article.Id, Categories: newCa}
			if err = InsertCategories(tx, newCa); err != nil {
				errCnt++
			}
			if err = a.InsertArticleCategory(tx); err != nil {
				errCnt++
			}
		}()

		// delete old categories
		go func() {
			defer wg.Done()
			a := Article{Id: article.Id, Categories: delCa}
			if err = a.DeleteArticleCategoryByBoth(tx); err != nil {
				errCnt++
			}
		}()
		wg.Wait()

		if errCnt != 0 {
			err = errors.New("error happened in UpdateArticleCmd()")
		}
		return
	})
	// End transaction

	return
}

// Insert draft.
type DraftCmd func(MySQL, Draft) error

func DraftCommand(mysql MySQL, draft Draft) (err error) {
	var d []Draft
	option := &service.QueryOption{
		Args: []interface{}{draft.ContentHash},
		Aom: map[string]service.Ope{
			"ContentHash": service.Eq,
		},
	}
	if d, err = FindDrafts(mysql.DB, option); err != nil {
		return
	}

	err = mysql.Transact(func(tx *sql.Tx) (err error) {
		if len(d) == 0 {
			service.GenNewId(service.IdLen, &draft.Id)
			if err = draft.InsertDraft(tx); err != nil {
				return
			}
			return
		}
		err = draft.UpdateDraft(tx)
		return
	})
	return
}

// Find articles.
type FindArticleCmd func(MySQL, *service.QueryOption) ([]Article, error)

func FindArticleCommand(mysql MySQL, option *service.QueryOption) ([]Article, error) {
	articles, err := FindArticle(mysql.DB, option)
	return articles, err
}

// Find articles by category id.
type FindArticleByCategoryCmd func(MySQL, []string, *service.QueryOption) ([]Article, error)

func FindArticleByCategoryCommand(mysql MySQL, categoryNames []string, option *service.QueryOption) ([]Article, error) {
	articles, err := FindArticleByCategoryId(mysql.DB, categoryNames, option)
	return articles, err
}

// Find categories.
type FindCategoryCmd func(MySQL, *service.QueryOption) ([]CategoryResponse, error)

func FindCategoryCommand(mysql MySQL, option *service.QueryOption) ([]CategoryResponse, error) {
	categories, err := FindCategory(mysql.DB, option)
	return categories, err
}

// Find drafts.
type FindDraftCmd func(MySQL, *service.QueryOption) ([]Draft, error)

func FindDraftCommand(mysql MySQL, option *service.QueryOption) ([]Draft, error) {
	drafts, err := FindDrafts(mysql.DB, option)
	return drafts, err
}

// Find the number of total articles.
type FindArticleNumCmd func(MySQL, *service.QueryOption) (int, error)

func FindArticlesNumCommand(mysql MySQL, option *service.QueryOption) (int, error) {
	articleNum, err := FindArticlesNum(mysql.DB, option)
	return articleNum, err
}

// Find the number of total articles by category id.
type FindArticlesNumByCategoryCmd func(MySQL, []string) (int, error)

func FindArticlesNumByCategoryCommand(mysql MySQL, categoryNames []string) (int, error) {
	articleNum, err := FindArticlesNumByCategoryId(mysql.DB, categoryNames)
	return articleNum, err
}

// Find the number of total drafts.
type FindDraftNumCmd func(MySQL, *service.QueryOption) (int, error)

func FindDraftsNumCommand(mysql MySQL, option *service.QueryOption) (int, error) {
	draftNum, err := FindDraftsNum(mysql.DB, option)
	return draftNum, err
}

// Insert category array to categories table.
func InsertCategories(tx *sql.Tx, categories []Category) (err error) {
	errCnt := 0
	wg := new(sync.WaitGroup)
	for _, c := range categories {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := c.InsertCategory(tx); err != nil {
				errCnt++
			}
		}()
	}
	wg.Wait()

	if errCnt != 0 {
		err = errors.New("error happened in InsertCategories()")
	}
	return
}
