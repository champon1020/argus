package database

import (
	"github.com/champon1020/argus"
	mgorm "github.com/champon1020/minigorm"
)

var (
	db *mgorm.DB
)

func init() {
	_db, err := mgorm.New("mysql", "root:toor@tcp(172.30.0.3:3306)/argus")
	if err != nil {
		argus.Logger.Fatalf("%v\n", err)
	}
	db = &_db
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
