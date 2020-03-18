package repo

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/champon1020/argus"
	"github.com/stretchr/testify/assert"
)

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
		WithArgs("TEST_ID").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow("TEST_CA_ID", "c1"))

	article := Article{
		Id: "TEST_ID",
		Categories: []Category{
			{"TEST_CA_ID2", "d1"},
			{"TEST_CA_ID3", "d2"},
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
	assert.Equal(t, Category{Id: "TEST_CA_ID2", Name: "d1"}, newCa[0])

	if len(delCa) == 0 {
		t.Fatalf("delCa is empty")
	}
	assert.Equal(t, Category{Id: "TEST_CA_ID", Name: "c1"}, delCa[0])
}

func TestExtractNewAndDelCategory(t *testing.T) {
	allCa := []Category{
		{"TEST_CA_ID2", "c2"},
		{"TEST_CA_ID3", "c3"},
	}
	existCa := []Category{
		{"TEST_CA_ID1", "c1"},
		{"TEST_CA_ID2", "c2"},
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
	assert.Equal(t, Category{Id: "TEST_CA_ID3", Name: "c3"}, newCa[0])

	if len(delCa) == 0 {
		t.Fatalf("delCa is empty")
	}
	assert.Equal(t, Category{Id: "TEST_CA_ID1", Name: "c1"}, delCa[0])
}
