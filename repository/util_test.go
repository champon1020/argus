package repository

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/champon1020/argus"
	"github.com/stretchr/testify/assert"
)

func TestConvertArticleId(t *testing.T) {
	db, mock, err := sqlmock.New()
	mysql := MySQL{}
	mysql.DB = db
	if err != nil {
		t.Fatalf("unable to create db mock")
	}
	defer db.Close()

	// ConvertArticleId()
	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT (id+1) FROM articles " +
			"WHERE (id+1) NOT IN " +
			"(SELECT id FROM articles) LIMIT ?")).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	article := Article{Id: -1}
	if err := ConvertArticleId(mysql, &article); err != nil {
		argus.StdLogger.ErrorLog(*Errors)
		t.Fatalf("error was occured in testing function\n")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("different from expectation")
	}
	assert.Equal(t, 1, article.Id)
}

func TestConvertCategoriesId(t *testing.T) {
	db, mock, err := sqlmock.New()
	mysql := MySQL{}
	mysql.DB = db
	if err != nil {
		t.Fatalf("unable to create db mock")
	}
	defer db.Close()

	// ConvertCategoriesId()
	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT (id+1) FROM categories " +
			"WHERE (id+1) NOT IN " +
			"(SELECT id FROM categories) LIMIT ?")).
		WithArgs(3).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(1).AddRow(3).AddRow(4))

	categories := []Category{
		{Id: -1},
		{Id: 2},
		{Id: -1},
	}
	if err := ConvertCategoriesId(mysql, &categories); err != nil {
		argus.StdLogger.ErrorLog(*Errors)
		t.Fatalf("error was occured in testing function\n")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("different from expectation")
	}
	assert.Equal(t, 1, categories[0].Id)
	assert.Equal(t, 2, categories[1].Id)
	assert.Equal(t, 3, categories[2].Id)
}

func TestConvertDraftId(t *testing.T) {
	db, mock, err := sqlmock.New()
	mysql := MySQL{}
	mysql.DB = db
	if err != nil {
		t.Fatalf("unable to create db mock")
	}
	defer db.Close()

	// ConvertDraftId()
	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT (id+1) FROM drafts " +
			"WHERE (id+1) NOT IN " +
			"(SELECT id FROM drafts) LIMIT ?")).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	draft := Draft{Id: -1}
	if err := ConvertDraftId(mysql, &draft); err != nil {
		argus.StdLogger.ErrorLog(*Errors)
		t.Fatalf("error was occured in testing function\n")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("different from expectation")
	}
	assert.Equal(t, 1, draft.Id)
}

func TestExtractCategory(t *testing.T) {
	db, mock, err := sqlmock.New()
	mysql := MySQL{}
	mysql.DB = db
	if err != nil {
		t.Fatalf("unable to create db mock")
	}
	defer db.Close()

	// FindCategoryByArticleId()
	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM categories " +
			"WHERE id IN (" +
			"SELECT category_id FROM article_category " +
			"WHERE article_id=?) ORDER BY name")).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(1, "c1"))

	article := Article{
		Id: 1,
		Categories: []Category{
			{2, "d1"},
			{3, "d2"},
		},
	}
	newCa, delCa, err := ExtractCategory(db, article)
	if err != nil {
		argus.StdLogger.ErrorLog(*Errors)
		t.Fatalf("error was occured in testing function\n")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("different from expectation")
	}

	if len(newCa) == 0 {
		t.Fatalf("newCa is empty")
	}
	assert.Equal(t, Category{Id: 2, Name: "d1"}, newCa[0])

	if len(delCa) == 0 {
		t.Fatalf("delCa is empty")
	}
	assert.Equal(t, Category{Id: 1, Name: "c1"}, delCa[0])
}

func TestExtractNewAndDelCategory(t *testing.T) {
	allCa := []Category{
		{2, "c2"},
		{3, "c3"},
	}
	existCa := []Category{
		{1, "c1"},
		{2, "c2"},
	}

	var (
		newCa, delCa []Category
		err          error
	)
	if newCa, delCa, err = ExtractNewAndDelCategory(allCa, existCa); err != nil {
		argus.StdLogger.ErrorLog(*Errors)
		t.Fatalf("error was occured in testing function\n")
	}
	if len(newCa) == 0 {
		t.Fatalf("newCa is empty")
	}
	assert.Equal(t, Category{Id: 3, Name: "c3"}, newCa[0])

	if len(delCa) == 0 {
		t.Fatalf("delCa is empty")
	}
	assert.Equal(t, Category{Id: 1, Name: "c1"}, delCa[0])
}
