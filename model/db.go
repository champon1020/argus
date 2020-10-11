package model

import (
	"errors"

	"github.com/champon1020/argus"
	"github.com/champon1020/minigorm"

	// Import mysql driver.
	_ "github.com/go-sql-driver/mysql"
)

// Db contains database or transaction instance.
var Db *Database

// InitDatabase initializes model.Database instance.
func InitDatabase() {
	Db = new(Database)
	Db.Connect(&argus.Config.Db)
}

var (
	errDbFailedConnect = errors.New("model.db: Failed to connect database")
)

// DatabaseIface is the interface of the Database struct.
type DatabaseIface interface {
	// Connect to database.
	Connect(config *argus.DbConf)

	// Article
	FindArticleByID(a *Article, id string) error
	FindAllArticles(a *[]Article, op *QueryOptions) error
	FindPublicArticles(a *[]Article, op *QueryOptions) error
	FindPublicArticlesGeSortedID(a *[]Article, sortedID int, op *QueryOptions) error
	FindPublicArticlesByTitle(a *[]Article, title string, op *QueryOptions) error
	FindPublicArticlesByCategory(a *[]Article, categoryID string, op *QueryOptions) error
	RegisterArticle(a *Article) error
	UpdateArticle(a *Article) error

	// Category
	FindPublicCategories(c *[]Category, op *QueryOptions) error

	// Draft
	FindDrafts(d *[]Draft, op *QueryOptions) error
	FindDraftByID(d *Draft, id string) error
	RegisterDraft(d *Draft) error
	UpdateDraft(d *Draft) error
	DeleteDraft(draftID string) error

	// Count
	CountAllArticles(cnt *int, op *QueryOptions) error
	CountPublicArticles(cnt *int, op *QueryOptions) error
	CountPublicArticlesByTitle(cnt *int, title string, op *QueryOptions) error
	CountPublicArticlesByCategory(cnt *int, categoryID string, op *QueryOptions) error
	CountDrafts(cnt *int, op *QueryOptions) error
}

// Database contains minigorm.DB.
type Database struct {
	DB *minigorm.DB
	TX *minigorm.TX
}

// Connect initializes database settings.
func (db *Database) Connect(config *argus.DbConf) {
	var err error
	dataSource :=
		config.User + ":" +
			config.Pass + "@tcp(" +
			config.Host + ":" +
			config.Port + ")/" +
			config.DbName + "?parseTime=true"

	if db.DB, err = minigorm.NewWithConf(minigorm.SourceConf{
		Driver:       "mysql",
		DataSource:   dataSource,
		MaxIdleConns: 50,
		MaxOpenConns: 100,
	}); err != nil {
		err = argus.NewError(errDbFailedConnect, err)
		argus.Logger.Fatalf("%v\n", err)
	}
}

// MockDatabase is the mock Database struct for test.
type MockDatabase struct{}

// Connect is dummy function.
// This function is declared for implementing DatabaseIface.
func (db *MockDatabase) Connect(config *argus.DbConf) {
	// dummy function
}

// QueryOptions is the struct includes options about sql query.
type QueryOptions struct {
	// limit query
	Limit int

	// offset query
	Offset int

	// orderby query
	OrderBy string

	// orderby direction (descending or ascending)
	Desc bool
}

// Apply the query options to context.
func (op *QueryOptions) apply(ctx *minigorm.Context) {
	if op == nil {
		return
	}

	if op.Limit > 0 {
		ctx.Limit(op.Limit)
	}

	if op.Offset > 0 {
		ctx.Offset(op.Offset)
	}

	if op.OrderBy != "" {
		if op.Desc {
			ctx.OrderByDesc(op.OrderBy)
		} else {
			ctx.OrderBy(op.OrderBy)
		}
	}
}

// NewOp create new QueryOptions object.
func NewOp(limit int, offset int, orderby string, desc bool) *QueryOptions {
	op := new(QueryOptions)

	if limit > 0 {
		op.Limit = limit
	}

	if offset > 0 {
		op.Offset = offset
	}

	if orderby != "" {
		op.OrderBy = orderby
	}

	if !desc {
		op.Desc = desc
	}

	return op
}
