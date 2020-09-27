package gosql

import (
	"database/sql"
	"fmt"
)

type Pool struct {
	list map[string]*sql.DB
}

func NewPool() *Pool {
	pool := &Pool{}
	pool.list = make(map[string]*sql.DB)
	return pool
}

func (p Pool) New(alias string, config *Config) error {
	var db *sql.DB
	var err error
	db, err = New(config)
	if nil != err {
		return fmt.Errorf("new db error occur %s", err)
	}
	// @todo 先写死
	// db.SetConnMaxLifetime(time.Minute * 5)
	// db.SetMaxIdleConns(5)
	// db.SetMaxOpenConns(20)
	p.SetClient(alias, db)
	return nil
}

func (p Pool) Close() error {
	for _, db := range p.GetPoolList() {
		db.Close()
	}
	return nil
}

func (p *Pool) GetClient(alias string) *sql.DB {
	if db, ok := p.list[alias]; ok {
		return db
	}
	return nil
}

func (p *Pool) GetPoolList() map[string]*sql.DB {
	return p.list
}

func (p *Pool) SetClient(alias string, client *sql.DB) error {
	if _, ok := p.list[alias]; ok {
		return fmt.Errorf("重复初始化db:%s", alias)
	}
	p.list[alias] = client
	return nil
}
