package repository

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/champon1020/argus"
	"github.com/champon1020/argus/service"
	"github.com/stretchr/testify/assert"
)

var (
	loc, _   = time.LoadLocation("Asia/Tokyo")
	testTime = time.Date(2020, 3, 9, 0, 0, 0, 0, loc)
)

func TestRegisterArticleCmd(t *testing.T) {
	db, mock, err := sqlmock.New()
	mysql := MySQL{}
	mysql.DB = db
	if err != nil {
		t.Fatalf("unable to create db mock")
	}
	defer db.Close()

	article := Article{
		Id:    -1,
		Title: "test",
		Categories: []Category{
			{Id: 1, Name: "c1"},
			{Id: 2, Name: "c2"},
		},
		CreateDate:  testTime,
		UpdateDate:  testTime,
		ContentHash: "0123456789",
		ImageHash:   "9876543210",
		Private:     false,
	}

	// GetMinId() with articles table
	// Called by repository/util.go: ConvertArticleId()
	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT (id+1) FROM articles " +
			"WHERE (id+1) NOT IN " +
			"(SELECT id FROM articles) LIMIT ?")).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	// FindCategoryByArticleId()
	// Called by repository/util.go: ExtractCategory()
	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM categories " +
			"WHERE id IN (" +
			"SELECT category_id FROM article_category " +
			"WHERE article_id=?) ORDER BY name")).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "c1"))

	// GetMinId() with categories table
	// Called by repository/util.go: ConvertCategoriesId()
	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT (id+1) FROM categories " +
			"WHERE (id+1) NOT IN " +
			"(SELECT id FROM categories) LIMIT ?")).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1).AddRow(2))

	// FindDrafts() with content hash
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM drafts WHERE content_hash=? ORDER BY id DESC OFFSET ?")).
		WithArgs("0123456789", 0).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	// Start transaction
	mock.ExpectBegin()

	// InsertCategories()
	mock.ExpectExec("INSERT INTO categories").
		WithArgs(2, "c2", "c2").
		WillReturnResult(sqlmock.NewResult(2, 1))

	// InsertArticles()
	mock.ExpectExec("INSERT INTO articles").
		WithArgs(1, "test", testTime, testTime, "0123456789", "9876543210", false).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// InsertArticleCategory()
	mock.ExpectExec("INSERT INTO article_category").
		WithArgs(1, 2).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Commit
	mock.ExpectCommit()

	if err := RegisterArticleCommand(mysql, article); err != nil {
		argus.StdLogger.ErrorLog(*Errors)
		t.Fatalf("error was occured in testing function\n")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("different from expectation: %v\n", err)
	}
}

func TestDraftCmd_Insert(t *testing.T) {
	db, mock, err := sqlmock.New()
	mysql := MySQL{}
	mysql.DB = db
	if err != nil {
		t.Fatalf("unable to create db mock")
	}
	defer db.Close()

	draft := Draft{
		Id:          -1,
		Title:       "draft",
		Categories:  "c1&c2",
		UpdateDate:  testTime,
		ContentHash: "0123456789",
		ImageHash:   "9876543210",
	}

	// FindDrafts()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM drafts WHERE content_hash=? ORDER BY id DESC OFFSET ?")).
		WithArgs("0123456789", 0).
		WillReturnRows(sqlmock.NewRows([]string{}))

	mock.ExpectBegin()

	// ConvertDraftId()
	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT (id+1) FROM drafts " +
			"WHERE (id+1) NOT IN " +
			"(SELECT id FROM drafts) LIMIT ?")).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	// InsertDraft()
	mock.ExpectExec("INSERT INTO drafts").
		WithArgs(1, "draft", "c1&c2", testTime, "0123456789", "9876543210").
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	if err := DraftCommand(mysql, draft); err != nil {
		argus.StdLogger.ErrorLog(*Errors)
		t.Fatalf("error was occured in testing function\n")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("different from expectation: %v\n", err)
	}
}

func TestDraftCmd_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	mysql := MySQL{}
	mysql.DB = db
	if err != nil {
		t.Fatalf("unable to create db mock")
	}
	defer db.Close()

	draft := Draft{
		Id:          1,
		Title:       "draft",
		Categories:  "c1&c2",
		UpdateDate:  testTime,
		ContentHash: "0123456789",
		ImageHash:   "9876543210",
	}

	// FindDrafts()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM drafts WHERE content_hash=? ORDER BY id DESC OFFSET ?")).
		WithArgs("0123456789", 0).
		WillReturnRows(
			sqlmock.NewRows([]string{
				"id", "title", "categories", "update_date", "content_hash", "image_hash",
			}).AddRow(1, "draft", "d1&d2", testTime, "0123456789", "9876543210"))

	mock.ExpectBegin()

	// UpdateDraft()
	mock.ExpectExec("UPDATE drafts").
		WithArgs("draft", "c1&c2", testTime, "0123456789", "9876543210", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	if err := DraftCommand(mysql, draft); err != nil {
		argus.StdLogger.ErrorLog(*Errors)
		t.Fatalf("error was occured in testing function\n")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("different from expectation: %v\n", err)
	}
}

func TestFindArticleCmd_All(t *testing.T) {
	db, mock, err := sqlmock.New()
	mysql := MySQL{}
	mysql.DB = db
	if err != nil {
		t.Fatalf("unable to create db mock")
	}
	defer db.Close()

	article := Article{}
	argFlg := service.GenFlg(article, "Limit")

	// FindArticle()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM articles ORDER BY id DESC LIMIT ? OFFSET ?")).
		WithArgs(0, 0).
		WillReturnRows(
			sqlmock.NewRows([]string{
				"id", "title", "create_date", "update_date", "content_hash", "image_hash", "private",
			}).AddRow(1, "test", testTime, testTime, "0123456789", "9876543210", false))

	// FindCategoriesByArticleId()
	// Called by article.go: FindArticle()
	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM categories " +
			"WHERE id IN (" +
			"SELECT category_id FROM article_category " +
			"WHERE article_id=?) ORDER BY name")).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "c1"))

	articles, err := FindArticleCommand(mysql, article, argFlg, 0)

	if err != nil {
		argus.StdLogger.ErrorLog(*Errors)
		t.Fatalf("error was occured in testing function\n")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("different from expectation: %v\n", err)
	}

	assert.Equal(t, len(articles), 1)
	assert.Equal(t, len(articles[0].Categories), 1)
	assert.Equal(t, articles[0].Categories[0].Name, "c1")
}

func TestFindArticleCmd_Title(t *testing.T) {
	db, mock, err := sqlmock.New()
	mysql := MySQL{}
	mysql.DB = db
	if err != nil {
		t.Fatalf("unable to create db mock")
	}
	defer db.Close()

	article := Article{Title: "test"}
	argFlg := service.GenFlg(article, "Title", "Limit")

	// FindArticle()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM articles WHERE title=? ORDER BY id DESC LIMIT ? OFFSET ?")).
		WithArgs("test", 0, 0).
		WillReturnRows(
			sqlmock.NewRows([]string{
				"id", "title", "create_date", "update_date", "content_hash", "image_hash", "private",
			}).AddRow(1, "test", testTime, testTime, "0123456789", "9876543210", false))

	// FindCategoriesByArticleId()
	// Called by article.go: FindArticle()
	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM categories " +
			"WHERE id IN (" +
			"SELECT category_id FROM article_category " +
			"WHERE article_id=?) ORDER BY name")).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "c1"))

	articles, err := FindArticleCommand(mysql, article, argFlg, 0)

	if err != nil {
		argus.StdLogger.ErrorLog(*Errors)
		t.Fatalf("error was occured in testing function\n")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("different from expectation: %v\n", err)
	}

	assert.Equal(t, len(articles), 1)
	assert.Equal(t, len(articles[0].Categories), 1)
	assert.Equal(t, articles[0].Categories[0].Name, "c1")
}

func TestFindArticleByCategoryCmd(t *testing.T) {
	db, mock, err := sqlmock.New()
	mysql := MySQL{}
	mysql.DB = db
	if err != nil {
		t.Fatalf("unable to create db mock")
	}
	defer db.Close()

	categoryNames := []string{"c1", "c2"}
	argFlg := service.GenFlg(Article{}, "Limit")

	// FindArticleByCategoryId()
	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM articles WHERE id IN ("+
			"SELECT article_id FROM article_category "+
			"WHERE category_id IN ("+
			"SELECT id FROM categories "+
			"WHERE name=? AND name=? )) LIMIT ? OFFSET ?")).
		WithArgs("c1", "c2", 0, 0).
		WillReturnRows(
			sqlmock.NewRows([]string{
				"id", "title", "create_date", "update_date", "content_hash", "image_hash", "private",
			}).AddRow(1, "test", testTime, testTime, "0123456789", "9876543210", false))

	// FindCategoriesByArticleId()
	// Called by article.go: FindArticle()
	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM categories " +
			"WHERE id IN (" +
			"SELECT category_id FROM article_category " +
			"WHERE article_id=?) ORDER BY name")).
		WithArgs(1).
		WillReturnRows(
			sqlmock.NewRows([]string{
				"id", "name",
			}).AddRow(1, "c1"))

	articles, err := FindArticleByCategoryCommand(mysql, categoryNames, argFlg, 0)

	if err != nil {
		argus.StdLogger.ErrorLog(*Errors)
		t.Fatalf("error was occured in testing function\n")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("different from expectation: %v\n", err)
	}

	assert.Equal(t, len(articles), 1)
	assert.Equal(t, len(articles[0].Categories), 1)
	assert.Equal(t, articles[0].Categories[0].Name, "c1")
}

func TestFindCategoryCmd(t *testing.T) {
	db, mock, err := sqlmock.New()
	mysql := MySQL{}
	mysql.DB = db
	if err != nil {
		t.Fatalf("unable to create db mock")
	}
	defer db.Close()

	category := Category{}

	// FindCategory()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM categories")).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(1, "c1"))

	// FindArticleNumByCategoryId()
	// Called by category.go: FindCategory()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(article_id) FROM article_category WHERE category_id=?")).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"articleNum"}).
			AddRow(3))

	categories, err := FindCategoryCommand(mysql, category, 0)

	if err != nil {
		argus.StdLogger.ErrorLog(*Errors)
		t.Fatalf("error was occured in testing function\n")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("different from expectation: %v\n", err)
	}

	assert.Equal(t, len(categories), 1)
	assert.Equal(t, categories[0].Name, "c1")
}

func TestFindDraftCmd(t *testing.T) {
	db, mock, err := sqlmock.New()
	mysql := MySQL{}
	mysql.DB = db
	if err != nil {
		t.Fatalf("unable to create db mock")
	}
	defer db.Close()

	draft := Draft{}
	argFlg := service.GenFlg(draft, "Limit")

	// FindDraft()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM drafts ORDER BY id DESC LIMIT ? OFFSET ?")).
		WithArgs(0, 0).
		WillReturnRows(
			sqlmock.NewRows([]string{
				"id", "title", "categories", "update_date", "content_hash", "image_hsah",
			}).AddRow(1, "draft", "c1&c2", testTime, "0123456789", "9876543210"))

	if _, err := FindDraftCommand(mysql, draft, argFlg, 0); err != nil {
		argus.StdLogger.ErrorLog(*Errors)
		t.Fatalf("error was occured in testing function\n")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("different from expectation: %v\n", err)
	}
}

func TestFindArticlesNumCommand(t *testing.T) {
	db, mock, err := sqlmock.New()
	mysql := MySQL{}
	mysql.DB = db
	if err != nil {
		t.Fatalf("unable to create db mock")
	}
	defer db.Close()

	article := Article{}
	argsFlg := service.GenFlg(article)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(id) FROM articles")).
		WillReturnRows(
			sqlmock.NewRows([]string{"articlesNum"}).AddRow(1))

	articlesNum, err := FindArticlesNumCommand(mysql, article, argsFlg)

	if err != nil {
		argus.StdLogger.ErrorLog(*Errors)
		t.Fatalf("error was occured in testing function\n")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("different from expectation: %v\n", err)
	}

	assert.Equal(t, articlesNum, 1)
}

func TestFindArticlesNumByCategoryCommand(t *testing.T) {
	db, mock, err := sqlmock.New()
	mysql := MySQL{}
	mysql.DB = db
	if err != nil {
		t.Fatalf("unable to create db mock")
	}
	defer db.Close()

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT COUNT(id) FROM articles " +
			"WHERE id IN (" +
			"SELECT article_id FROM article_category " +
			"WHERE category_id IN (" +
			"SELECT id FROM categories " +
			"WHERE name=? ))")).
		WithArgs("c1").
		WillReturnRows(
			sqlmock.NewRows([]string{"articlesNum"}).AddRow(1))

	articlesNum, err := FindArticlesNumByCategoryCommand(mysql, []string{"c1"})

	if err != nil {
		argus.StdLogger.ErrorLog(*Errors)
		t.Fatalf("error was occured in testing function\n")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("different from expectation: %v\n", err)
	}

	assert.Equal(t, articlesNum, 1)
}

func TestFindDraftsNumCommand(t *testing.T) {
	db, mock, err := sqlmock.New()
	mysql := MySQL{}
	mysql.DB = db
	if err != nil {
		t.Fatalf("unable to create db mock")
	}
	defer db.Close()

	draft := Draft{}
	argsFlg := service.GenFlg(draft)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(id) FROM drafts")).
		WillReturnRows(
			sqlmock.NewRows([]string{"draftsNum"}).AddRow(1))

	draftsNum, err := FindDraftsNumCommand(mysql, draft, argsFlg)

	if err != nil {
		argus.StdLogger.ErrorLog(*Errors)
		t.Fatalf("error was occured in testing function\n")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("different from expectation: %v\n", err)
	}

	assert.Equal(t, draftsNum, 1)
}
