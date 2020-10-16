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
		ReplicationConfig: *NewReplicationConfig(&config.Replication),
		PartitionConfig:   *NewPartitionConfig(&config.Partition),
	}
}

type ClusterConfig struct {
	ReplicationConfig ReplicationConfig
	PartitionConfig   PartitionConfig
}

func (c *ClusterConfig) GetReplicationConfig() *ReplicationConfig {
	return &c.ReplicationConfig
}

func (c *ClusterConfig) GetPartitionConfig() *PartitionConfig {
	return &c.PartitionConfig
}

func NewCluster(conf Configer, exec func(*sql.DB, DBHandler, TableSelector)) (*Cluster, error) {

	// 分库设置
	dp, err := NewPartition(conf.GetClusterConfig().GetPartitionConfig())

	if !errors.Is(err, NopDbPartitionError) && conf.GetClusterConfig().GetPartitionConfig().GetDbTotal() != len(conf.GetShardsConfig()) {
		return nil, fmt.Errorf("cluster must have %d shards, but config %d", len(conf.GetShardsConfig()), conf.GetClusterConfig().GetPartitionConfig().GetDbTotal())
	}

	var shards []*Shard
	for i, cs := range conf.GetShardsConfig() {
		// 设置分库
		cs.SetDbname(dp.Pick(cs.GetDbname(), i))

		shards = append(shards, NewShard(cs, AddExecutor(exec), AddTablePartition(dp)))
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

func (p *Cluster) Select(key []byte, table string) *Shard {
	return p.Shards[p.DP.Partition(key, table)]
}

func (p *Cluster) Pick() *Shard {
	return p.Shards[0]
}

func (p *Cluster) Map(f func(*Shard)) {
	for _, shard := range p.Shards {
		f(shard)
	}
}
