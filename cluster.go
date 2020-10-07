package gosql

import (
	"database/sql"
	"errors"
	"fmt"
)

func NewConfig(conf *YamlConfig) *Config {

	var shards []*ShardConfig
	for _, shard := range conf.Shards {
		shards = append(shards, NewShardConfig(shard, &conf.Db, &conf.Cluster.Replication))
	}

	return &Config{
		Cluster: NewClusterConfig(&conf.Cluster),
		Shards:  shards,
	}
}

type Config struct {
	Cluster *ClusterConfig
	Shards  []*ShardConfig
}

func (c *Config) GetClusterConfig() *ClusterConfig {
	return c.Cluster
}

func (c *Config) GetShardsConfig() []*ShardConfig {
	return c.Shards
}

func NewClusterConfig(config *YamlClusterConfig) *ClusterConfig {
	return &ClusterConfig{
		ReplicationConfig:    *NewReplicationConfig(&config.Replication),
		DbPartitionConfig:    *NewDbPartitionConfig(&config.Partition),
		TablePartitionConfig: *NewTablePartitionConfig(&config.Partition),
	}
}

type ClusterConfig struct {
	ReplicationConfig    ReplicationConfig
	DbPartitionConfig    DbPartitionConfig
	TablePartitionConfig TablePartitionConfig
}

func (c *ClusterConfig) GetReplicationConfig() *ReplicationConfig {
	return &c.ReplicationConfig
}

func (c *ClusterConfig) GetDbPartitionConfig() *DbPartitionConfig {
	return &c.DbPartitionConfig
}

func (c *ClusterConfig) GetTablePartitionConfig() *TablePartitionConfig {
	return &c.TablePartitionConfig
}

func NewCluster(conf Configer, exec func(*sql.DB, DBHandler, TablePartition)) (*Cluster, error) {

	// 分库设置
	dp, err := NewDbPartition(conf.GetClusterConfig().GetDbPartitionConfig())

	if !errors.Is(err, NopDbPartitionError) && conf.GetClusterConfig().GetDbPartitionConfig().GetTotal() != len(conf.GetShardsConfig()) {
		return nil, fmt.Errorf("cluster must have %d shards, but config %d", len(conf.GetShardsConfig()), conf.GetClusterConfig().GetDbPartitionConfig().GetTotal())
	}

	// 分表
	tp := NewTablePartition(conf.GetClusterConfig().GetTablePartitionConfig())
	var shards []*Shard
	for i, cs := range conf.GetShardsConfig() {
		// 设置分库
		cs.SetDbname(dp.Divide(cs.GetDbname(), i))
		shards = append(shards, NewShard(cs, AddExecutor(exec), AddTablePartition(tp)))
	}

	return &Cluster{
		Shards: shards,
		DP:     dp,
	}, nil

}

type Cluster struct {
	DP     DbPartitioner
	Shards []*Shard
}

func (p *Cluster) Write(handler DBHandler) {
	p.Pick().Write(handler)
}

func (p *Cluster) Read(handler DBHandler) {
	p.Pick().Read(handler)
}

func (p *Cluster) Select(key []byte) *Shard {
	return p.Shards[p.DP.Select(key)]
}

func (p *Cluster) Pick() *Shard {
	return p.Shards[0]
}

func (p *Cluster) Map(f func(*Shard)) {
	for _, shard := range p.Shards {
		f(shard)
	}
}
