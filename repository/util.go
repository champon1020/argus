package repository

import (
	"database/sql"

	"github.com/champon1020/argus/service"
)

func GenArgsQuery(argsFlg uint32, st interface{}) string {
	return service.GenArgsQuery(argsFlg, st)
}

func GenArgsSlice(argsFlg uint32, st interface{}) []interface{} {
	return service.GenArgsSliceLogic(argsFlg, st, false)
}

func GenArgsSliceIsLimit(argsFlg uint32, st interface{}, isLimit bool) []interface{} {
	return service.GenArgsSliceLogic(argsFlg, st, isLimit)
}

// Get and Set empty and minimum article id.
func ArticleIdConverter(mysql MySQL, article *Article) (err error) {
	var idList []int
	if idList, err = GetEmptyMinId(mysql.DB, "articles", 1); err != nil {
		return
	}

	(*article).Id = idList[0]
	return
}

// Get and Set category empty minimum id.
func CategoriesIdConverter(mysql MySQL, categories *[]Category) (err error) {
	var idList []int
	if idList, err = GetEmptyMinId(mysql.DB, "categories", len(*categories)); err != nil {
		return
	}

	cur := 0
	for i := 0; i < len(*categories); i++ {
		if (*categories)[i].Id == -1 {
			(*categories)[i].Id = idList[cur]
			cur++
		}
	}
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

func GetEmptyMinId(db *sql.DB, tableName string, numOfId int) (res []int, err error) {
	query := "SELECT (id+1) FROM " +
		tableName + " " +
		"WHERE (id+1) NOT IN " +
		"(SELECT id FROM " +
		tableName + " " +
		") LIMIT ?"

	rows, err := db.Query(query, numOfId)
	defer func() {
		if err = rows.Close(); err != nil {
			logger.ErrorPrintf(err)
			return
		}
	}()

	if err != nil {
		logger.ErrorPrintf(err)
		return
	}

	if rows == nil {
		logger.ErrorPrintf(err)
		return
	}

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			logger.ErrorPrintf(err)
			break
		}
		res = append(res, id)
	}
	return
}
