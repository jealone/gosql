package gosql

import (
	"bytes"
	"fmt"
)

type Sharding interface {
	Select(string, []byte, ...interface{}) (int, int, *bytes.Buffer)
	Allocation(int, string, ...interface{}) string
	GetDbname() string
	TableSelector() TableSelector
}

type StandardSharding struct {
	dbname   string
	conf     ShardingConfig
	shards   int
	selector TableSelector
}

func (s *StandardSharding) Select(table string, key []byte, params ...interface{}) (int, int, *bytes.Buffer) {
	ti, total, t := s.selector.Select(table, key, params...)
	return ti / (total / s.shards), ti, t
}

func (s *StandardSharding) Allocation(i int, dbname string, params ...interface{}) string {
	b := AcquireBuffer()
	defer ReleaseBuffer(b)
	b.WriteString(dbname)
	b.WriteByte('_')
	b.WriteString(fmt.Sprintf("%02x", i))
	return b.String()
}

func (s *StandardSharding) TableSelector() TableSelector {
	return s.selector
}

func (s *StandardSharding) GetDbname() string {
	return s.dbname
}

func ProvideStandardSharding(conf *Config, selector TableSelector) *StandardSharding {
	c := conf.GetShardingConfig()

	return &StandardSharding{
		conf:     c,
		dbname:   c.GetDbname(),
		shards:   c.GetTotal(),
		selector: selector,
	}
}
