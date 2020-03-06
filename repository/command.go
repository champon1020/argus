package repository

import (
	"database/sql"
	"sync"

	"github.com/champon1020/argus/service"
)

/*
Flow:
	- Convert article id.
	- Insert article.
	- Sample new categories from article.Categories.
	- Insert new categories.
	- Insert the pair of article_id and category_ids.
*/
func RegisterArticleCmd(mysql MySQL, article Article) (err error) {
	if err = ArticleIdConverter(mysql, &article); err != nil {
		return
	}

	var newCa []Category
	if newCa, _, err = ExtractCategory(mysql.DB, article); err != nil {
		return
	}
	if err = CategoriesIdConverter(mysql, &newCa); err != nil {
		return
	}
	article.Categories = newCa

	// Start transaction
	err = mysql.Transact(func(tx *sql.Tx) (err error) {
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
	return
}

/*
Flow:
	- Get categories by article_id.
	- Classify categories to new and old.
	- Update article.
	- Insert new categories.
	- Insert the pair of article_id and new category_ids.
	- Delete the pair of article_id and old category_ids.
*/
func UpdateArticleCmd(mysql MySQL, article Article) (err error) {
	var newCa, delCa []Category
	if newCa, delCa, err = ExtractCategory(mysql.DB, article); err != nil {
		return
	}
	if err = CategoriesIdConverter(mysql, &newCa); err != nil {
		return
	}

	// Start transaction
	err = mysql.Transact(func(tx *sql.Tx) (err error) {
		wg := new(sync.WaitGroup)
		wg.Add(3)

		// update article
		go func() {
			defer wg.Done()
			if err = article.UpdateArticle(tx); err != nil {
				return
			}
		}()

		// insert new categories
		go func() {
			defer wg.Done()
			if err = CategoriesIdConverter(mysql, &newCa); err != nil {
				return
			}
			a := Article{Id: article.Id, Categories: newCa}
			if err = InsertCategories(tx, newCa); err != nil {
				return
			}
			if err = a.InsertArticleCategory(tx); err != nil {
				return
			}
		}()

		// delete old categories
		go func() {
			defer wg.Done()
			a := Article{Id: article.Id, Categories: delCa}
			if err = a.DeleteArticleCategoryByBoth(tx); err != nil {
				return
			}
		}()
		wg.Wait()

		return
	})
	return
}

func FindArticleCmd(mysql MySQL, article Article, argFlg uint32) (articles []Article, err error) {
	articles, err = article.FindArticle(mysql.DB, argFlg)
	return
}

// Get articles by category.
func FindArticleByCategoryCmd(mysql MySQL, categoryNames []string) (articles []Article, err error) {
	articles, err = FindArticleByCategoryId(mysql.DB, categoryNames)
	return
}

func FindCategoryCmd(mysql MySQL, category Category, argFlg uint32) (categories []CategoryResponse, err error) {
	categories, err = category.FindCategory(mysql.DB, argFlg)
	return
}

// Insert category array to categories table.
func InsertCategories(tx *sql.Tx, categories []Category) (err error) {
	wg := new(sync.WaitGroup)
	for _, c := range categories {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err = c.InsertCategory(tx); err != nil {
				return
			}
		}()
	}
	wg.Wait()
	return
}

// Extract new, exist, or deleted category
// from category array found by article_id from article_category table.
// - newCa: categories which are added to inserted or updated article
// - delCa: categories which are removed from inserted or updated article
func ExtractCategory(db *sql.DB, article Article) (newCa, delCa []Category, err error) {
	var existCa, bufCa []Category
	if existCa, err = article.FindCategoryByArticleId(db); err != nil {
		return
	}
	if bufCa, delCa, err = ExtractNewAndDelCategory(article.Categories, existCa); err != nil {
		return
	}
	for _, c := range bufCa {
		var ca []CategoryResponse
		if ca, err = c.FindCategory(db, service.GenFlg(Category{}, "Name")); err != nil {
			return
		}
		if len(ca) == 0 {
			continue
		}
		newCa = append(newCa, Category{Id: ca[0].Id, Name: ca[0].Name})
	}
	return
}

// Extract new, del category.
func ExtractNewAndDelCategory(allCa, existCa []Category) (newCa, delCa []Category, err error) {
	cMap := map[string]Category{}
	for _, c := range existCa {
		cMap[c.Name] = c
	}

	for i := 0; i < len(allCa); i++ {
		if _, ok := cMap[allCa[i].Name]; !ok {
			newCa = append(newCa, allCa[i])
			continue
		}
		delete(cMap, allCa[i].Name)
	}

	for _, c := range cMap {
		delCa = append(delCa, c)
	}
	return
}
