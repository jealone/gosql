package gosql

import (
	"fmt"
	"net/url"
	"time"
)

func ProvideConfig(config *YamlConfig) (*Config, error) {

	var ss []ShardingTableConfig

	for _, c := range config.Cluster.Sharding.Table {
		if 0 != c.Total%config.Cluster.Sharding.DbTotal {
			return nil, fmt.Errorf("table total must be a multiple of db total")
		}
		ss = append(ss, ShardingTableConfig{
			Name:  c.Name,
			Total: c.Total,
		})
	}

	var scs []ShardConfig
	for _, shard := range config.Shards {
		sc := ShardConfig{
			MasterConfig: DbConfig{
				Driver:   config.Cluster.Driver,
				Dbname:   config.Cluster.Dbname,
				Host:     shard.Master.Host,
				User:     shard.Master.User,
				Password: shard.Master.Password,
				Protocol: shard.Master.Protocol,
				Params:   shard.Master.Params,
				Conn: ConnConfig{
					MaxLifeTime:  config.Cluster.Conn.MaxLifetime,
					MaxIdleConns: config.Cluster.Conn.MaxIdleConns,
					MaxOpenConns: config.Cluster.Conn.MaxOpenConns,
				},
			},
		}
		for _, replica := range shard.Replicas.Hosts {
			sc.ReplicasConfig = append(sc.ReplicasConfig, DbConfig{
				Driver:   config.Cluster.Driver,
				Dbname:   config.Cluster.Dbname,
				Host:     replica,
				User:     shard.Replicas.User,
				Password: shard.Replicas.Password,
				Protocol: shard.Replicas.Protocol,
				Params:   shard.Replicas.Params,
				Conn: ConnConfig{
					MaxLifeTime:  config.Cluster.Conn.MaxLifetime,
					MaxIdleConns: config.Cluster.Conn.MaxIdleConns,
					MaxOpenConns: config.Cluster.Conn.MaxOpenConns,
				},
			})
		}
		scs = append(scs, sc)
	}

	if len(scs) != config.Cluster.Sharding.DbTotal {
		return nil, fmt.Errorf("Shards total unequal db instance total")
	}

	return &Config{
		ShardsConfig: scs,
		ShardingConfig: ShardingConfig{
			Dbname: config.Cluster.Dbname,
			Type:   config.Cluster.Type,
			Total:  config.Cluster.Sharding.DbTotal,
			Table:  ss,
		}}, nil
}

type Config struct {
	ShardingConfig ShardingConfig
	ShardsConfig   []ShardConfig
}

func (c Config) GetShardsConfig() []ShardConfig {
	return c.ShardsConfig
}

func (c Config) GetShardingConfig() ShardingConfig {
	return c.ShardingConfig
}

type ShardingConfig struct {
	Dbname string
	Total  int
	Type   string
	Table  []ShardingTableConfig
}

func (s ShardingConfig) GetDbname() string {
	return s.Dbname
}

func (s ShardingConfig) GetType() string {
	return s.Type
}

func (s ShardingConfig) GetTotal() int {
	return s.Total
}

func (s ShardingConfig) GetTableConfig() []ShardingTableConfig {
	return s.Table
}

type ShardingTableConfig struct {
	Name  string
	Total int
}

func (s *ShardingTableConfig) GetTableName() string {
	return s.Name
}

func (s *ShardingTableConfig) GetTotal() int {
	return s.Total
}

type ShardConfig struct {
	MasterConfig   DbConfig
	ReplicasConfig []DbConfig
}

func (c ShardConfig) GetMasterConfig() DbConfig {
	return c.MasterConfig
}

func (c ShardConfig) GetReplicasConfig() []DbConfig {
	return c.ReplicasConfig
}

type DbConfig struct {
	Driver   string
	Dbname   string
	Protocol string
	Host     string
	User     string
	Password string
	Params   map[string]string
	Conn     ConnConfig
}

func (c DbConfig) GetDriver() string {
	return c.Driver
}

func (c DbConfig) GetDbname() string {
	return c.Dbname
}

func (c DbConfig) GetProtocol() string {
	return c.Protocol
}

func (c DbConfig) GetHost() string {
	return c.Host
}

func (c DbConfig) GetUser() string {
	return c.User
}

func (c DbConfig) GetPassword() string {
	return c.Password
}

func (c DbConfig) GetParams() map[string]string {
	return c.Params
}

func (c DbConfig) GetConn() ConnConfig {
	return c.Conn
}

type ConnConfig struct {
	MaxLifeTime  int
	MaxIdleConns int
	MaxOpenConns int
}

func (c ConnConfig) GetMaxLifeTime() time.Duration {
	return time.Duration(c.MaxLifeTime) * time.Millisecond
}

func (c ConnConfig) GetMaxIdleConns() int {
	return c.MaxIdleConns
}

func (c ConnConfig) GetMaxOpenConns() int {
	return c.MaxOpenConns
}

func (c DbConfig) GetUrl(dbname string) string {
	b := AcquireBuffer()
	defer ReleaseBuffer(b)

	b.WriteString(c.User)
	b.WriteByte(':')
	b.WriteString(c.Password)

	b.WriteByte('@')
	if len(c.Protocol) > 0 {
		b.WriteString(c.Protocol)
	} else {
		b.WriteString("tcp")
	}

	b.WriteByte('(')
	b.WriteString(c.Host)
	b.WriteByte(')')

	b.WriteByte('/')

	if 0 == len(dbname) {
		b.WriteString(c.Dbname)
	} else {
		b.WriteString(dbname)
	}

	if total := len(c.Params); total > 0 {
		params := make(url.Values, total)
		for key, v := range c.Params {
			params.Set(key, v)
		}
		if p := params.Encode(); 0 != len(p) {
			b.WriteByte('?')
			b.WriteString(p)
		}
	}

	return b.String()
}
