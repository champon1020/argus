package repository

import (
	"database/sql"
	"sync"
)

func RegisterArticle(mysql MySQL, article Article) {
	mysql.Transact(func(tx *sql.Tx) (err error) {
		wg := new(sync.WaitGroup)
		wg.Add(2)

		go func() {
			defer wg.Done()
			err = article.InsertArticle(tx)
		}()

		go func() {
			defer wg.Done()
			err = article.InsertArticleCategory(tx)
		}()

		innerWg := new(sync.WaitGroup)
		for _, c := range article.Categories {
			innerWg.Add(1)
			go func() {
				defer innerWg.Done()

				err = c.InsertCategory(tx)
				if err != nil {
					logger.ErrorPrintf(err)
				}
			}()
		}
		innerWg.Wait()
		wg.Wait()

		return
	})
}

func UpdateArticle(mysql MySQL, article Article) {
	mysql.Transact(func(tx *sql.Tx) (err error) {
		nowCategories := article.FindArticleCategory(mysql.DB)
		cMap := map[int]Category{}
		for _, c := range nowCategories {
			cMap[c.Id] = c
		}

		var newCategories, delCategories []Category
		for i := 0; i < len(article.Categories); i++ {
			if c, ok := cMap[article.Categories[i].Id]; !ok {
				newCategories = append(newCategories, c)
				continue
			}
			delete(cMap, article.Categories[i].Id)
		}

		for _, c := range cMap {
			delCategories = append(delCategories, c)
		}

		wg := new(sync.WaitGroup)
		wg.Add(2)

		// insert new categories
		go func() {
			defer wg.Done()

			a := Article{Id: article.Id, Categories: newCategories}
			err = a.InsertArticleCategory(tx)
			if err != nil {
				logger.ErrorPrintf(err)
			}
		}()

		// delete old categories
		go func() {
			defer wg.Done()

			a := Article{Id: article.Id, Categories: delCategories}
			err = a.DeleteArticleCategoryByBoth(tx)
			if err != nil {
				logger.ErrorPrintf(err)
			}
		}()
		wg.Wait()

		return err
	})
}

func FindArticle(mysql MySQL, article Article, argFlg uint32) (articles []Article) {
	articles = article.FindArticle(mysql.DB, argFlg)
	return
}

func FindCategory(mysql MySQL, category Category, argFlg uint32) (categories []Category) {
	categories = category.FindCategory(mysql.DB, argFlg)
	return
}
