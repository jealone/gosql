package gosql

import (
	"bytes"
	"database/sql"
	"fmt"
)

type DBHandler func(*sql.DB, int, *bytes.Buffer) error

var DefaultDBExecutor DbExecutor = func(db *sql.DB, handler DBHandler, node int, table *bytes.Buffer) {
	_ = handler(db, node, table)
}

type StaticShard struct {
	Master  *sql.DB
	Replica []*sql.DB
	lb      Replication
}

func (s *StaticShard) GetMaster() *sql.DB {
	return s.Master
}

func (s *StaticShard) GetReplica() *sql.DB {
	total := s.GetReplicaTotal()
	if 0 == total {
		return s.Master
	}
	i := s.lb.Replicate(total)
	if i < total {
		return s.Replica[i]
	} else {
		return s.Replica[i%total]
	}
}

func (s *StaticShard) GetReplicaTotal() int {
	return len(s.Replica)
}

func ProvideStaticShards(conf *Config, sharding Sharding, lb Replication) ([]Shard, error) {
	if "static" != conf.GetShardingConfig().GetType() {
		return nil, fmt.Errorf("static cluster type must be static ")
	}
	var shards []Shard
	for i, c := range conf.GetShardsConfig() {
		if 1 == conf.GetShardingConfig().GetTotal() {
			shards = append(shards, NewShard(c, sharding.GetDbname(), lb))
		} else {
			shards = append(shards, NewShard(c, sharding.Allocation(i, sharding.GetDbname()), lb))
		}
	}
	return shards, nil
}

func NewShard(conf ShardConfig, dbname string, lb Replication) *StaticShard {

	conf.GetMasterConfig()
	master := NewDB(conf.GetMasterConfig().GetDriver(), conf.GetMasterConfig().GetUrl(dbname), conf.GetMasterConfig().GetConn())

	var replicas []*sql.DB

	for _, c := range conf.GetReplicasConfig() {
		replicas = append(replicas, NewDB(c.GetDriver(), c.GetUrl(dbname), c.GetConn()))
	}

	return &StaticShard{
		Master:  master,
		Replica: replicas,
		lb:      lb,
	}
}
