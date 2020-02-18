package repository

import (
	"database/sql"

	"github.com/champon1020/argus"
	_ "github.com/go-sql-driver/mysql"
)

var (
	logger argus.Logger
	config argus.Config
)

func init() {
	logger.NewLogger("[repository]")
	config.Load()
}

type MySQL struct {
	*sql.DB
}

func (mysql *MySQL) Connect(config argus.DbConf, dbName string) (err error) {
	dataSourceName :=
		config.User + ":" +
			config.Pass + "@tcp(" +
			config.Host + ":" +
			config.Port + ")/" +
			dbName + "?parseTime=true"

	logger.Printf("DataSource: %s\n", dataSourceName)
	if mysql.DB, err = sql.Open("mysql", dataSourceName); err != nil {
		logger.ErrorPrintf(err)
		return
	}

	mysql.SetMaxIdleConns(20)
	mysql.SetConnMaxLifetime(1)
	mysql.SetMaxOpenConns(10)

	return nil
}

func (mysql *MySQL) Transact(txFunc func(*sql.Tx) error) (err error) {
	tx, err := mysql.Begin()
	if err != nil {
		logger.ErrorPrintf(err)
		return
	}

	err = txFunc(tx)
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			logger.Printf("%v\n", p)
			return
		} else if err != nil {
			tx.Rollback()
			logger.ErrorPrintf(err)
			return
		} else {
			err = tx.Commit()
		}
	}()

	return
}
