package repo

import (
	"database/sql"
)

// Extract new, exist, or deleted category
// from category array found by article_id from article_category table.
// - newCa: categories which are added to inserted or updated articles
// - delCa: categories which are removed from inserted or updated articles
func ExtractCategory(db *sql.DB, article Article) (newCa, delCa []Category, err error) {
	var existCa []Category
	if existCa, err = article.FindCategoryByArticleId(db); err != nil {
		return
	}
	if newCa, delCa, err = ExtractNewAndDelCategory(article.Categories, existCa); err != nil {
		return
	}
	return
}

// Extract new, del category.
func ExtractNewAndDelCategory(allCa, existCa []Category) (newCa, delCa []Category, err error) {
	cMap := make(map[string]Category)
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
