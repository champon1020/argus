package repository

import (
	"database/sql"

	"github.com/champon1020/argus"
)

var logger argus.Logger

func init() {
	logger.NewLogger("[Repository]")
}

type MySQL struct {
	*sql.DB
}

func (mysql *MySQL) Connect(config argus.Config, dbName string) (err error) {
	dataSourceName :=
		config.Db.User + ":" +
			config.Db.Pass + "@tcp(" +
			config.Db.Host + ":" +
			config.Db.Port + ")/" +
			dbName
	mysql.DB, err = sql.Open("mysql", dataSourceName)
	mysql.SetMaxIdleConns(5)
	mysql.SetConnMaxLifetime(1)
	mysql.SetMaxOpenConns(10)
	return
}

func (mysql *MySQL) Transact(
	txFunc func(*sql.Tx) error) (err error) {

	tx, err := mysql.Begin()
	if err != nil {
		logger.ErrorPrintf(err)
		return
	}

	err = txFunc(tx)
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			logger.ErrorPanic(err)
		} else if err != nil {
			tx.Rollback()
			logger.ErrorPanic(err)
		} else {
			err = tx.Commit()
		}
	}()

	return
}
