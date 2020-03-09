package repository

import (
	"regexp"
	"testing"
	"time"

	"github.com/champon1020/argus"

	"github.com/DATA-DOG/go-sqlmock"
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
	mock.ExpectQuery(regexp.QuoteMeta("SELECT (id+1) FROM articles WHERE (id+1) NOT IN (SELECT id FROM articles ) LIMIT ?")).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	// FindCategoryByArticleId()
	// Called by repository/util.go: ExtractCategory()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM categories WHERE id IN (SELECT category_id FROM article_category WHERE article_id=?) ORDER BY name")).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "c1"))

	// FindCategory() with category name
	// Called by repository/util.go: ExtractCategory()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM categories WHERE name=?")).
		WithArgs("c2").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(2, "c2"))

	// FindArticleNumByCategoryId()
	// Called by category.go: FindCategory()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(article_id) FROM article_category WHERE category_id=?")).
		WithArgs(2).
		WillReturnRows(sqlmock.NewRows([]string{"articleNum"}).AddRow(1))

	// GetMinId() with categories table
	// Called by repository/util.go: ConvertCategoriesId()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT (id+1) FROM categories WHERE (id+1) NOT IN (SELECT id FROM categories ) LIMIT ?")).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1).AddRow(2))

	// FindDrafts() with content hash
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM drafts WHERE content_hash=? ORDER BY id DESC LIMIT ?")).
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

	if err := RegisterArticleCmd(mysql, article); err != nil {
		argus.StdLogger.ErrorLog(*Errors)
		t.Errorf("error was occured in testing function\n")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("different from expectation")
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
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM drafts WHERE content_hash=? ORDER BY id DESC LIMIT ?")).
		WithArgs("0123456789", 0).
		WillReturnRows(sqlmock.NewRows([]string{}))

	mock.ExpectBegin()

	// ConvertDraftId()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT (id+1) FROM drafts WHERE (id+1) NOT IN (SELECT id FROM drafts ) LIMIT ?")).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	// InsertDraft()
	mock.ExpectExec("INSERT INTO drafts").
		WithArgs(1, "draft", "c1&c2", testTime, "0123456789", "9876543210").
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	if err := DraftCmd(mysql, draft); err != nil {
		argus.StdLogger.ErrorLog(*Errors)
		t.Errorf("error was occured in testing function\n")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("different from expectation")
	}
}
