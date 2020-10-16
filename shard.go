package gosql

import (
	"database/sql"
	"net/url"
	"time"
)

func NewShardConfig(shard *YamlShardConfig, db *YamlDBConfig, r *YamlReplicationConfig) *ShardConfig {
	return &ShardConfig{
		YamlShardConfig: *shard,
		YamlDBConfig:    *db,
		tr:              *NewReplicationConfig(r),
	}
}

type ShardConfig struct {
	YamlDBConfig
	YamlShardConfig
	tr ReplicationConfig
}

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

type DBHandler func(*sql.DB, TableSelector) error

type Option interface {
	apply(*Shard)
}

type optionFunc func(*Shard)

func (f optionFunc) apply(c *Shard) {
	f(c)
}

func NewShard(conf *ShardConfig, options ...Option) *Shard {

	c := &Shard{
		Executor: exec,
		Master:   NewMasterNode(conf.GetDriver(), conf),
		Replica:  NewReplicasNode(conf.GetDriver(), conf, AddReplication(NewReplication(&conf.tr))),
		selector: GetNopTablePartition(),
	}

	return c.WithOptions(options...)
}

func AddExecutor(f func(*sql.DB, DBHandler, TableSelector)) Option {
	return optionFunc(func(s *Shard) {
		if nil != f {
			s.Executor = f
		}
	})
}

func AddTablePartition(selector TableSelector) Option {
	return optionFunc(func(s *Shard) {
		if nil != selector {
			s.selector = selector
		} else {
			s.selector = GetNopTablePartition()
		}
	})
}

func AddMaster(conf ShardConfiger) Option {
	return optionFunc(func(s *Shard) {
		s.Master = NewMasterNode(conf.GetDriver(), conf)
	})
}

func AddReplica(conf ShardConfiger, opts ...ReplicaNodeOption) Option {
	return optionFunc(func(s *Shard) {
		s.Replica = NewReplicasNode(conf.GetDriver(), conf, opts...)
	})
}

type Shard struct {
	Master   *MasterNode
	Replica  *ReplicaNode
	selector TableSelector
	Executor func(*sql.DB, DBHandler, TableSelector)
}

func (s *Shard) clone() *Shard {
	copied := *s
	return &copied
}

func (s *Shard) WithOptions(opts ...Option) *Shard {
	shard := s.clone()
	for _, opt := range opts {
		opt.apply(shard)
	}
	return shard
}

func exec(db *sql.DB, handler DBHandler, selector TableSelector) {
	_ = handler(db, selector)
}

func (s *Shard) Write(handler DBHandler) {
	s.Executor(s.Master.GetDB(), handler, s.selector)
}

func (s *Shard) Read(handler DBHandler) {
	db := s.Replica.GetDB()

	if nil != db {
		db = s.Master.GetDB()
	}

	s.Executor(db, handler, s.selector)
}
