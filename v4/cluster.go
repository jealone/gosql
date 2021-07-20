package gosql

import "database/sql"

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
