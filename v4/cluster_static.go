package gosql

import (
	"bytes"
	"database/sql"
)

type DbExecutor func(*sql.DB, DBHandler, int, *bytes.Buffer)

type StaticCluster struct {
	sharding Sharding
	shards   []Shard
	config   *Config
	executor DbExecutor
}

func (s *StaticCluster) Write(table string, key []byte, handler DBHandler, params ...interface{}) {
	di, ti, t := s.sharding.Select(table, key, params...)
	defer ReleaseBuffer(t)
	s.executor(s.shards[di].GetMaster(), handler, ti, t)
}

func (s *StaticCluster) Read(table string, key []byte, handler DBHandler, params ...interface{}) {
	di, ti, t := s.sharding.Select(table, key, params...)
	defer ReleaseBuffer(t)
	s.executor(s.shards[di].GetReplica(), handler, ti, t)
}

func (s *StaticCluster) TableSelector() TableSelector {
	return s.sharding.TableSelector()
}

func (s *StaticCluster) GetShardsTotal() int {
	return len(s.shards)
}
