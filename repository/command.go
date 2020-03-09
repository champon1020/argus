package repository

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
	flg := service.GenFlg(Draft{}, "ContentHash")
	draft := Draft{ContentHash: article.ContentHash}
	if d, err = draft.FindDrafts(mysql.DB, flg); err != nil {
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

type DraftCmd func(MySQL, Draft) error

func DraftCommand(mysql MySQL, draft Draft) (err error) {
	var d []Draft
	flg := service.GenFlg(Draft{}, "ContentHash")
	if d, err = draft.FindDrafts(mysql.DB, flg); err != nil {
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

type FindArticleCmd func(MySQL, Article, uint32) ([]Article, error)

func FindArticleCommand(mysql MySQL, article Article, argFlg uint32) (articles []Article, err error) {
	articles, err = article.FindArticle(mysql.DB, argFlg)
	return
}

type FindArticleByCategoryCmd func(MySQL, []string, uint32) ([]Article, error)

// Get articles by category.
func FindArticleByCategoryCommand(mysql MySQL, categoryNames []string, argFlg uint32) (articles []Article, err error) {
	articles, err = FindArticleByCategoryId(mysql.DB, categoryNames, argFlg)
	return
}

type FindCategoryCmd func(MySQL, Category, uint32) ([]CategoryResponse, error)

func FindCategoryCommand(mysql MySQL, category Category, argFlg uint32) (categories []CategoryResponse, err error) {
	categories, err = category.FindCategory(mysql.DB, argFlg)
	return
}

type FindDraftCmd func(MySQL, Draft, uint32) ([]Draft, error)

func FindDraftCommand(mysql MySQL, draft Draft, argFlg uint32) (drafts []Draft, err error) {
	drafts, err = draft.FindDrafts(mysql.DB, argFlg)
	return
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
