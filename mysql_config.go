package qlap

import (
	"database/sql"

	"github.com/go-sql-driver/mysql"
)

type MysqlConfig struct {
	*mysql.Config
}

func (myCfg *MysqlConfig) openAndPing(maxIdleConns int) (*sql.DB, error) {
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
