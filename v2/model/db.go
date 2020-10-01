package model

import (
	"errors"

	"github.com/champon1020/argus/v2"
	mgorm "github.com/champon1020/minigorm"
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
	Connect(config *argus.DbConf)
}

// Database contains mgorm.DB.
type Database struct {
	DB *mgorm.DB
	TX *mgorm.TX
}

// Connect of Database initializes database settings.
func (db *Database) Connect(config *argus.DbConf) {
	dataSource :=
		config.User + ":" +
			config.Pass + "@tcp(" +
			config.Host + ":" +
			config.Port + ")/" +
			config.DbName + "?parseTime=true"
	_db, err := mgorm.New("mysql", dataSource)
	if err != nil {
		err = argus.NewError(errDbFailedConnect, err)
		argus.Logger.Fatalf("%v\n", err)
	}
	db.DB = &_db
}

// MockDatabase is the mock Database struct for test.
type MockDatabase struct{}

// Connect of MockDatabase is dummy function.
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
func (op *QueryOptions) apply(ctx *mgorm.Context) {
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
