package repository

import (
	"database/sql"
	"sync"
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
	if err = articleIdConverter(mysql, &article); err != nil {
		return
	}

	if err = categoriesIdConverter(mysql, &article.Categories); err != nil {
		return
	}

	err = mysql.Transact(func(tx *sql.Tx) (err error) {
		wg := new(sync.WaitGroup)
		wg.Add(3)

		go func() {
			defer wg.Done()
			if err = article.InsertArticle(tx); err != nil {
				return
			}
		}()

		go func() {
			defer wg.Done()
			if err = InsertCategories(tx, article.Categories); err != nil {
				return
			}
		}()

		go func() {
			defer wg.Done()
			if err = article.InsertArticleCategory(tx); err != nil {
				return
			}
		}()

		wg.Wait()
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
	err = mysql.Transact(func(tx *sql.Tx) (err error) {
		nowCategories, err := article.FindArticleCategory(mysql.DB)
		if err != nil {
			return
		}

		cMap := map[int]Category{}
		for _, c := range nowCategories {
			cMap[c.Id] = c
		}

		var newCategories, delCategories []Category
		for i := 0; i < len(article.Categories); i++ {
			if _, ok := cMap[article.Categories[i].Id]; !ok {
				newCategories = append(newCategories, article.Categories[i])
				continue
			}
			delete(cMap, article.Categories[i].Id)
		}

		for _, c := range cMap {
			delCategories = append(delCategories, c)
		}

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
			if len(newCategories) > 0 {
				if err = categoriesIdConverter(mysql, &newCategories); err != nil {
					return
				}
				a := Article{Id: article.Id, Categories: newCategories}
				if err = InsertCategories(tx, newCategories); err != nil {
					return
				}
				if err = a.InsertArticleCategory(tx); err != nil {
					return
				}
			}
		}()

		// delete old categories
		go func() {
			defer wg.Done()
			if len(delCategories) > 0 {
				a := Article{Id: article.Id, Categories: delCategories}
				if err = a.DeleteArticleCategoryByBoth(tx); err != nil {
					return
				}
			}
		}()
		wg.Wait()

		return err
	})
	return
}

func FindArticleCmd(mysql MySQL, article Article, argFlg uint32) (articles []Article, err error) {
	articles, err = article.FindArticle(mysql.DB, argFlg)
	return
}

func FindCategoryCmd(mysql MySQL, category Category, argFlg uint32) (categories []Category, err error) {
	categories, err = category.FindCategory(mysql.DB, argFlg)
	return
}

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
