package gosql

import (
	"time"
)

type Configer interface {
	GetClusterConfig() *ClusterConfig
	GetShardsConfig() []*ShardConfig
}

type ClusterConfiger interface {
	GetReplicationConfig() *ReplicationConfig
	GetDbPartitionConfig() *DbPartitionConfig
	GetTablePartitionConfig() *TablePartitionConfig
}

type DbPartitionConfiger interface {
	GetTotal() int
	GetType() string
}

type TablePartitionConfiger interface {
	GetType() string
	GetSpecs() []*TablePartitionSpecsConfig
}

type TablePartitionSpecConfiger interface {
	GetName() string
	GetTotal() int
}

type ReplicationConfiger interface {
	GetType() string
}

type DbConfiger interface {
	GetDriver() string
	GetDbname() string
	SetDbname(string)
	ConnConfiger
}

type ShardConfiger interface {
	MasterNodeConfiger
	ReplicaNodeConfiger
	ReplicationConfiger
}

type ConnConfiger interface {
	GetMaxLifetime() time.Duration
	GetMaxIdleConns() int
	GetMaxOpenConns() int
}

type MasterNodeConfiger interface {
	GetMasterUrl() string
	DbConfiger
}

type ReplicaNodeConfiger interface {
	GetReplicaUrls() []string
	DbConfiger
}

/*
func NewShardConfig(shard *YamlShardConfig, db *YamlDBConfig) *ShardConfig {
	if nil == shard || nil == db {
		panic("missing config")
	}

	return &ShardConfig{
		YamlShardConfig: *shard,
		YamlDBConfig:    *db,
	}
}

type ShardConfig struct {
	YamlShardConfig
	YamlDBConfig
}

// [user[:password]@][net[(addr)]]/dbname[?param1=value1&paramN=valueN]
func (s *ShardConfig) GetMasterUrl() string {

	b := AcquireBuffer()
	defer ReleaseBuffer(b)

	b.WriteString(s.Master.User)
	b.WriteByte(':')
	b.WriteString(s.Master.Password)

	b.WriteByte('@')
	if len(s.Master.Protocol) > 0 {
		b.WriteString(s.Master.Protocol)
	} else {
		b.WriteString("tcp")
	}

	b.WriteByte('(')
	b.WriteString(s.Master.Host)
	b.WriteByte(')')

	b.WriteByte('/')
	b.WriteString(s.GetDbname())

	params := make(url.Values)
	for key, v := range s.Master.Params {
		params.Set(key, v)
	}

	if p := params.Encode(); 0 != len(p) {
		b.WriteByte('?')
		b.WriteString(p)
	}

	return b.String()
}

func (s *ShardConfig) GetReplicaUrls() []string {
	var urls []string

	params := make(url.Values)
	for key, v := range s.Replicas.Params {
		params.Set(key, v)
	}

	paramStr := params.Encode()

	b := AcquireBuffer()
	defer ReleaseBuffer(b)

	for _, host := range s.Replicas.Hosts {

		b.WriteString(s.Replicas.User)
		b.WriteByte(':')
		b.WriteString(s.Replicas.Password)

		b.WriteByte('@')
		if len(s.Replicas.Protocol) > 0 {
			b.WriteString(s.Replicas.Protocol)
		} else {
			b.WriteString("tcp")
		}

		b.WriteByte('(')
		b.WriteString(host)
		b.WriteByte(')')

		b.WriteByte('/')
		b.WriteString(s.GetDbname())

		params := make(url.Values)
		for key, v := range s.Master.Params {
			params.Set(key, v)
		}

		if 0 != len(paramStr) {
			b.WriteByte('?')
			b.WriteString(paramStr)
		}

		urls = append(urls, b.String())

		b.Reset()
	}

	return urls
}

func (s *ShardConfig) GetDriver() string {
	return s.Driver
}

func (s *ShardConfig) GetDbname() string {
	return s.Dbname
}

func (s *ShardConfig) SetDbname(name string) {
	s.Dbname = name
}

func (s *ShardConfig) GetMaxLifetime() time.Duration {
	return time.Duration(s.Conn.MaxLifetime) * time.Millisecond
}

func (s *ShardConfig) GetMaxIdleConns() int {
	return s.Conn.MaxIdleConns
}

func (s *ShardConfig) GetMaxOpenConns() int {
	return s.Conn.MaxOpenConns
}

func NewConfig(config *YamlConfig) *Config {
	if nil == config {
		panic("missing config")
	}

	clusterConfig := NewClusterConfig(config.Cluster)

	if clusterConfig.Partition.Total != len(config.Shards) {
		panic("cluster db partition config error")
	}

	partitionConf := NewPartitionConfig(config.Cluster.Partition)

	var shardConfigs []*ShardConfig

	for i, shard := range config.Shards {

		c := NewShardConfig(shard, config.Db)

		partitionConf.Divide(c, i)

		shardConfigs = append(shardConfigs, c)
	}

	return &Config{
		Cluster: clusterConfig,
		Shards:  shardConfigs,
	}

}

type Config struct {
	Cluster *ClusterConfig
	Shards  []*ShardConfig
}

func NewClusterConfig(conf *YamlClusterConfig) *ClusterConfig {
	return &ClusterConfig{
		*conf,
	}
}

type ClusterConfig struct {
	YamlClusterConfig
}

func (c *ClusterConfig) GetTotal() int {
	return c.Partition.Total
}

func (c *ClusterConfig) GetType() string {
	return c.Partition.Type
}

func (c *Config) GetShardsConfig() interface{} {
	if 1 == len(c.Shards) {
		return c.Shards[0]
	}
	return c.Shards
}

func (c *Config) GetClusterConfig() ClusterConfiger {
	return c.Cluster
}

func NewPartitionConfig(conf *YamlPartitionConfig) *PartitionConfig {
	return &PartitionConfig{
		*conf,
	}
}

type PartitionConfig struct {
	YamlPartitionConfig
}

func (p *PartitionConfig) Divide(conf ShardConfiger, num int) {

	if p.Total > 1 {
		conf.SetDbname(DividesGenerate(conf.GetDbname(), num))
	}
}

*/
