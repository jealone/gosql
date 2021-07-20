package gosql

import (
	"bytes"
	"database/sql"
)

type DbExecutor func(*sql.DB, DBHandler, int, *bytes.Buffer)

type Cluster interface {
	Write(table string, key []byte, handler DBHandler, params ...interface{})
	Read(table string, key []byte, handler DBHandler, params ...interface{})
	GetShardsTotal() int
	TableSelector() TableSelector
}

type Shard interface {
	GetMaster() *sql.DB
	GetReplica() *sql.DB
	GetReplicaTotal() int
}
