package repository

import (
	"database/sql"
	"log"

	"github.com/champon1020/argus"
	_ "github.com/go-sql-driver/mysql"
)

var (
	GlobalMysql  MySQL
	Logger       = argus.Logger
	Errors       = &argus.Errors
	RuntimeError = argus.Error{Type: argus.DbRuntimeError}
	CmdError     = argus.Error{Type: argus.DbCmdFailedError}
	ScanError    = argus.Error{Type: argus.DbScanFailedError}
	QueryError   = argus.Error{Type: argus.DbQueryFailedError}
	CloseError   = argus.Error{Type: argus.DbCloseFailedError}
)

type MySQL struct {
	*sql.DB
}

func NewMysql() {
	if err := GlobalMysql.Connect(argus.GlobalConfig.Db); err != nil {
		log.Fatalf("%v\n", err)
	}
}

func (mysql *MySQL) Connect(config argus.DbConf) (err error) {
	dataSourceName :=
		config.User + ":" +
			config.Pass + "@tcp(" +
			config.Host + ":" +
			config.Port + ")/" +
			config.DbName + "?parseTime=true"
	Logger.Printf("DataSource: %s\n", dataSourceName)

	if mysql.DB, err = sql.Open("mysql", dataSourceName); err != nil {
		RuntimeError.SetErr(err).AppendTo(Errors)
		return
	}
	mysql.SetMaxIdleConns(20)
	mysql.SetConnMaxLifetime(1)
	mysql.SetMaxOpenConns(10)
	return
}

func (mysql *MySQL) Transact(txFunc func(*sql.Tx) error) (err error) {
	var tx *sql.Tx
	if tx, err = mysql.Begin(); err != nil {
		RuntimeError.SetErr(err).AppendTo(Errors)
		return
	}

	err = txFunc(tx)
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			Logger.Printf("%v\n", p)
			return
		} else if err != nil {
			tx.Rollback()
			Logger.Printf("%v", err)
			return
		} else {
			err = tx.Commit()
		}
	}()
	return
}

func RowsClose(rows *sql.Rows) {
	if rows == nil {
		return
	}
	if err := rows.Close(); err != nil {
		CloseError.SetErr(err).AppendTo(Errors)
	}
}
