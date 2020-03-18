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
	if err = ConvertArticleId(mysql, &article); err != nil {
		return
	}
	if newCa, _, err = ExtractCategory(mysql.DB, article); err != nil {
		return
	}
	if err = ConvertCategoriesId(mysql, &newCa); err != nil {
		return
	}
	article.Categories = newCa

	var d []Draft
	flg := service.GenMask(Draft{}, "ContentHash")
	draft := Draft{ContentHash: article.ContentHash}
	if d, err = draft.FindDrafts(mysql.DB, flg, [2]int{}); err != nil {
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
	if err = ConvertCategoriesId(mysql, &newCa); err != nil {
		return
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
			if err = ConvertCategoriesId(mysql, &newCa); err != nil {
				errCnt++
			}
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
	flg := service.GenMask(Draft{}, "ContentHash")
	if d, err = draft.FindDrafts(mysql.DB, flg, [2]int{}); err != nil {
		return
	}

	err = mysql.Transact(func(tx *sql.Tx) (err error) {
		if len(d) == 0 {
			if err = ConvertDraftId(mysql, &draft); err != nil {
				return
			}
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
type FindArticleCmd func(MySQL, Article, uint32, OffsetLimit) ([]Article, error)

func FindArticleCommand(mysql MySQL, article Article, argsMask uint32, lo OffsetLimit) ([]Article, error) {
	articles, err := article.FindArticle(mysql.DB, argsMask, lo)
	return articles, err
}

// Find articles by category id.
type FindArticleByCategoryCmd func(MySQL, []string, uint32, OffsetLimit) ([]Article, error)

func FindArticleByCategoryCommand(mysql MySQL, categoryNames []string, argsMask uint32, ol OffsetLimit) ([]Article, error) {
	articles, err := FindArticleByCategoryId(mysql.DB, categoryNames, argsMask, ol)
	return articles, err
}

// Find categories.
type FindCategoryCmd func(MySQL, Category, uint32, OffsetLimit) ([]CategoryResponse, error)

func FindCategoryCommand(mysql MySQL, category Category, argsMask uint32, ol OffsetLimit) ([]CategoryResponse, error) {
	categories, err := category.FindCategory(mysql.DB, argsMask, ol)
	return categories, err
}

// Find drafts.
type FindDraftCmd func(MySQL, Draft, uint32, OffsetLimit) ([]Draft, error)

func FindDraftCommand(mysql MySQL, draft Draft, argsMask uint32, ol OffsetLimit) ([]Draft, error) {
	drafts, err := draft.FindDrafts(mysql.DB, argsMask, ol)
	return drafts, err
}

// Find the number of total articles.
type FindArticleNumCmd func(MySQL, Article, uint32) (int, error)

func FindArticlesNumCommand(mysql MySQL, article Article, argsMask uint32) (int, error) {
	articleNum, err := article.FindArticlesNum(mysql.DB, argsMask)
	return articleNum, err
}

// Find the number of total articles by category id.
type FindArticlesNumByCategoryCmd func(MySQL, []string) (int, error)

func FindArticlesNumByCategoryCommand(mysql MySQL, categoryNames []string) (int, error) {
	articleNum, err := FindArticlesNumByCategoryId(mysql.DB, categoryNames)
	return articleNum, err
}

// Find the number of total drafts.
type FindDraftNumCmd func(MySQL, Draft, uint32) (int, error)

func FindDraftsNumCommand(mysql MySQL, draft Draft, argsMask uint32) (int, error) {
	draftNum, err := draft.FindDraftsNum(mysql.DB, argsMask)
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