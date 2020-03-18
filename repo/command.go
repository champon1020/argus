package repo

import (
	"database/sql"
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
			if err = d[0].DeleteDraft(tx); err != nil {
				return
			}
		}

		var (
			categoryId        string
			articleCategories []Category
		)
		wg := new(sync.WaitGroup)
		for _, c := range article.Categories {
			wg.Add(1)
			go func() {
				defer wg.Done()
				if categoryId, err = c.Exist(tx, &service.QueryOption{
					Args: []interface{}{c.Name},
					Aom:  map[string]service.Ope{"Name": service.Eq},
				}); err != nil {
					return
				}
				if categoryId == "" {
					service.GenNewId(service.IdLen, &c.Id)
					if err = c.InsertCategory(tx); err != nil {
						return
					}
				} else {
					c.Id = categoryId
				}
				articleCategories = append(articleCategories, c)
			}()
		}
		wg.Wait()

		// Insert into articles
		article.Categories = articleCategories
		if err = article.InsertArticle(tx); err != nil {
			return
		}

		// Insert into article_category
		ac := ArticleCategory{ArticleId: article.Id}
		wg = new(sync.WaitGroup)
		for _, c := range article.Categories {
			wg.Add(1)
			go func() {
				defer wg.Done()
				ac.CategoryId = c.Id
				if err = ac.InsertArticleCategory(tx); err != nil {
					return
				}
			}()
		}
		wg.Wait()
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
	// Start transaction
	err = mysql.Transact(func(tx *sql.Tx) (err error) {
		var (
			categoryId        string
			articleCategories []Category
		)
		wg := new(sync.WaitGroup)
		for _, c := range article.Categories {
			wg.Add(1)
			go func() {
				defer wg.Done()
				if categoryId, err = c.Exist(tx, &service.QueryOption{
					Args: []interface{}{c.Name},
					Aom:  map[string]service.Ope{"Name": service.Eq},
				}); err != nil {
					return
				}
				if categoryId == "" {
					service.GenNewId(service.IdLen, &c.Id)
					if err = c.InsertCategory(tx); err != nil {
						return
					}
				} else {
					c.Id = categoryId
				}
				articleCategories = append(articleCategories, c)
			}()
		}
		wg.Wait()
		article.Categories = articleCategories
		// delete func

		// Insert into articles
		if err = article.UpdateArticle(tx); err != nil {
			return
		}

		// Insert into article_category
		ac := ArticleCategory{ArticleId: article.Id}
		wg = new(sync.WaitGroup)
		for _, c := range article.Categories {
			wg.Add(1)
			go func() {
				defer wg.Done()
				ac.CategoryId = c.Id
				if err = ac.InsertArticleCategory(tx); err != nil {
					return
				}
			}()
		}
		wg.Wait()
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
