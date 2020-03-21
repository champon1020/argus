package repo

import (
	"database/sql"
	"strconv"
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
	args := &service.QueryArgs{
		Value: interface{}(article.ContentHash),
		Name:  "ContentHash",
		Ope:   service.Eq,
	}
	option := &service.QueryOption{
		Args: []*service.QueryArgs{args},
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
				args := &service.QueryArgs{
					Value: interface{}(c.Name),
					Name:  "Name",
					Ope:   service.Eq,
				}

				if categoryId, err = c.Exist(tx, &service.QueryOption{
					Args: []*service.QueryArgs{args},
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

				args := &service.QueryArgs{
					Value: interface{}(c.Name),
					Name:  "Name",
					Ope:   service.Eq,
				}

				if categoryId, err = c.Exist(tx, &service.QueryOption{
					Args: []*service.QueryArgs{args},
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

		// Insert into articles
		if err = article.UpdateArticle(tx); err != nil {
			return
		}

		// Delete from article_category
		args := &service.QueryArgs{
			Value: interface{}(article.Id),
			Name:  "ArticleId",
			Ope:   service.Eq,
		}
		option := &service.QueryOption{
			Args: []*service.QueryArgs{args},
		}
		for i, c := range article.Categories {
			option.Args = append(option.Args, &service.QueryArgs{
				Value: c.Id,
				Name:  "CategoryId#" + strconv.Itoa(i),
				Ope:   service.Ne,
			})
		}
		if err = DeleteArticleCategory(tx, option); err != nil {
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
	args := &service.QueryArgs{
		Value: interface{}(draft.ContentHash),
		Name:  "ContentHash",
		Ope:   service.Eq,
	}
	option := &service.QueryOption{
		Args: []*service.QueryArgs{args},
	}
	if d, err = FindDrafts(mysql.DB, option); err != nil {
		return
	}

	err = mysql.Transact(func(tx *sql.Tx) (err error) {
		if len(d) == 0 {
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

type DeleteDraftCmd func(MySQL, Draft) error

func DeleteDraftCommand(mysql MySQL, draft Draft) (err error) {
	err = mysql.Transact(func(tx *sql.Tx) (err error) {
		err = draft.DeleteDraft(tx)
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
type FindArticleByCategoryCmd func(MySQL, *service.QueryOption) ([]Article, error)

func FindArticleByCategoryCommand(mysql MySQL, option *service.QueryOption) ([]Article, error) {
	articles, err := FindArticleByCategoryId(mysql.DB, option)
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
type FindArticlesNumByCategoryCmd func(MySQL, *service.QueryOption) (int, error)

func FindArticlesNumByCategoryCommand(mysql MySQL, option *service.QueryOption) (int, error) {
	articleNum, err := FindArticlesNumByCategoryId(mysql.DB, option)
	return articleNum, err
}

// Find the number of total drafts.
type FindDraftNumCmd func(MySQL, *service.QueryOption) (int, error)

func FindDraftsNumCommand(mysql MySQL, option *service.QueryOption) (int, error) {
	draftNum, err := FindDraftsNum(mysql.DB, option)
	return draftNum, err
}
