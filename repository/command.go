package repository

import (
	"database/sql"
	"sync"
)

func RegisterArticleCmd(mysql MySQL, article Article) (err error) {
	if err = articleIdConverter(mysql, &article); err != nil {
		return
	}

	err = mysql.Transact(func(tx *sql.Tx) (err error) {
		wg := new(sync.WaitGroup)
		wg.Add(3)

		go func() {
			defer wg.Done()
			err = article.InsertArticle(tx)
			if err != nil {
				logger.ErrorPrintf(err)
				return
			}
		}()

		go func() {
			defer wg.Done()
			err = article.InsertArticleCategory(tx)
			if err != nil {
				logger.ErrorPrintf(err)
				return
			}
		}()

		go func() {
			defer wg.Done()
			innerWg := new(sync.WaitGroup)
			for _, c := range article.Categories {
				innerWg.Add(1)
				go func() {
					defer innerWg.Done()

					err = c.InsertCategory(tx)
					if err != nil {
						logger.ErrorPrintf(err)
						return
					}
				}()
			}
			innerWg.Wait()
		}()

		wg.Wait()
		return
	})
	return
}

func UpdateArticleCmd(mysql MySQL, article Article) (err error) {
	err = mysql.Transact(func(tx *sql.Tx) (err error) {
		nowCategories, err := article.FindArticleCategory(mysql.DB)
		if err != nil {
			logger.ErrorPrintf(err)
			return
		}

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
				return
			}
		}()

		// delete old categories
		go func() {
			defer wg.Done()

			a := Article{Id: article.Id, Categories: delCategories}
			err = a.DeleteArticleCategoryByBoth(tx)
			if err != nil {
				logger.ErrorPrintf(err)
				return
			}
		}()
		wg.Wait()

		return err
	})
	return
}

func FindArticleCmd(mysql MySQL, article Article, argFlg uint32) (articles []Article, err error) {
	articles, err = article.FindArticle(mysql.DB, argFlg)
	if err != nil {
		logger.ErrorPrintf(err)
	}
	return
}

func FindCategoryCmd(mysql MySQL, category Category, argFlg uint32) (categories []Category, err error) {
	categories, err = category.FindCategory(mysql.DB, argFlg)
	if err != nil {
		logger.ErrorPrintf(err)
	}
	return
}
