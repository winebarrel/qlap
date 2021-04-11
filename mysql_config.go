package qlap

import (
	"database/sql"

	"github.com/go-sql-driver/mysql"
)

type MysqlConfig struct {
	*mysql.Config
	OnlyPrint bool
}

type DB interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Close() error
}

func (myCfg *MysqlConfig) openAndPing(maxIdleConns int) (DB, error) {
	if myCfg.OnlyPrint {
		return &NullDB{}, nil
	}

	dsn := myCfg.FormatDSN()
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(0)
	db.SetMaxIdleConns(0)

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(maxIdleConns)

	return db, nil
}
