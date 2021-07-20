package gosql

import (
	"database/sql"
)

func NewDB(driver string, url string, c ConnConfig) *sql.DB {

	db, err := sql.Open(driver, url)
	if nil != err {
		panic(err)
	}

	db.SetConnMaxLifetime(c.GetMaxLifeTime())
	db.SetMaxIdleConns(c.GetMaxIdleConns())
	db.SetMaxOpenConns(c.GetMaxOpenConns())

	return db
}
