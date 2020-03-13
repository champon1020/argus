package repository

import (
	"database/sql"

	"github.com/champon1020/argus"
	_ "github.com/go-sql-driver/mysql"
)

var (
	GlobalMysql  *MySQL
	Errors       = &argus.Errors
	RuntimeError = argus.NewError(argus.DbRuntimeError)
	CmdError     = argus.NewError(argus.DbCmdFailedError)
	ScanError    = argus.NewError(argus.DbScanFailedError)
	QueryError   = argus.NewError(argus.DbQueryFailedError)
	CloseError   = argus.NewError(argus.DbCloseFailedError)
)

type MySQL struct {
	*sql.DB
}

func NewMysql() *MySQL {
	mysql := new(MySQL)
	if err := mysql.Connect(argus.GlobalConfig.Db); err != nil {
		argus.StdLogger.Fatalf("%v\n", err)
	}
	return mysql
}

func (mysql *MySQL) Connect(config argus.DbConf) (err error) {
	dataSourceName :=
		config.User + ":" +
			config.Pass + "@tcp(" +
			config.Host + ":" +
			config.Port + ")/" +
			config.DbName + "?parseTime=true"
	argus.Logger.Printf("DataSource: %s\n", dataSourceName)

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
			argus.Logger.Printf("%v\n", p)
			return
		} else if err != nil {
			tx.Rollback()
			argus.Logger.Printf("%v", err)
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
