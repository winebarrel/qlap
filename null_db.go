package qlap

import (
	"database/sql"
	"fmt"
	"os"
)

type NullDB struct{}

func (db *NullDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	fmt.Fprintln(os.Stderr, query)
	return nil, nil
}

func (db *NullDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	fmt.Fprintln(os.Stderr, query)
	return &sql.Rows{}, nil
}

func (db *NullDB) QueryRow(query string, args ...interface{}) *sql.Row {
	fmt.Fprintln(os.Stderr, query)
	return &sql.Row{}
}

func (db *NullDB) Close() error {
	return nil
}
