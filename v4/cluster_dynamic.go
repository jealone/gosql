package gosql

import (
	"bytes"
	"database/sql"
	"fmt"
	"sync"
	"time"
)

func ProvideDynamicRe(conf *Config) (map[string]*sql.DB, error) {
	if "dynamic" != conf.GetShardingConfig().GetType() {
		return nil, fmt.Errorf("dynamic cluster type must be dynamic ")
	}

	m := make(map[string]*sql.DB, conf.GetShardingConfig().GetTotal())
	return m, nil
}

type DynamicCluster struct {
	mu       sync.RWMutex `wire:"-"`
	Re       map[string]*sql.DB
	Config   *Config
	Sharding Sharding
	Selector DbSelector
	Executor DbExecutor
	Lb       Replication
}

func (s *DynamicCluster) Write(table string, key []byte, handler DBHandler, params ...interface{}) {
	di, ti, t := s.Sharding.Select(table, key, params...)
	d := s.Selector.Select(di, s.Sharding.GetDbname(), params...)

	defer ReleaseBuffer(d)

	conf := s.Config.GetShardsConfig()[di].GetMasterConfig()
	dbStr := conf.GetUrl(d.String())
	s.mu.RLock()
	if db, ok := s.Re[dbStr]; ok {
		s.mu.RUnlock()
		s.Executor(db, handler, ti, t)
		return
	}
	s.mu.RUnlock()

	s.mu.Lock()

	db := NewDB(conf.GetDriver(), dbStr, conf.GetConn())
	s.Re[dbStr] = db
	s.mu.Unlock()

	s.Executor(db, handler, ti, t)

	return
}

func (s *DynamicCluster) Read(table string, key []byte, handler DBHandler, params ...interface{}) {
	di, ti, t := s.Sharding.Select(table, key, params...)
	d := s.Selector.Select(di, s.Sharding.GetDbname(), params...)

	defer ReleaseBuffer(d)
	confs := s.Config.GetShardsConfig()[di].GetReplicasConfig()

	i := s.Lb.Replicate(len(confs))

	conf := confs[i]

	dbStr := conf.GetUrl(d.String())
	s.mu.RLock()
	if db, ok := s.Re[dbStr]; ok {
		s.mu.RUnlock()
		s.Executor(db, handler, ti, t)
		return
	}
	s.mu.RUnlock()

	s.mu.Lock()

	db := NewDB(conf.GetDriver(), dbStr, conf.GetConn())
	s.Re[dbStr] = db
	s.mu.Unlock()

	s.Executor(db, handler, ti, t)

	return
}

func (s *DynamicCluster) TableSelector() TableSelector {
	return s.Sharding.TableSelector()
}

func (s *DynamicCluster) GetShardsTotal() int {
	return len(s.Config.ShardsConfig)
}

type DbSelector interface {
	Select(int, string, ...interface{}) *bytes.Buffer
}

type DailyDbSelector struct {
}

func (se *DailyDbSelector) Select(di int, dbname string, params ...interface{}) *bytes.Buffer {

	b := AcquireBuffer()
	b.WriteString(dbname)
	b.WriteByte('_')
	b.WriteString(fmt.Sprintf("%02x", di))

	if len(params) < 1 {
		return b
	}

	if ts, ok := params[0].(int64); ok {
		b.WriteByte('_')
		b.WriteString(time.Unix(ts, 0).Format("20060102"))
	}

	return b

}

func ProvideDailyDbSelector() *DailyDbSelector {
	return &DailyDbSelector{}
}

type MonthlyDbSelector struct {
}

func (se *MonthlyDbSelector) Select(di int, dbname string, params ...interface{}) *bytes.Buffer {

	b := AcquireBuffer()
	b.WriteString(dbname)
	b.WriteByte('_')
	b.WriteString(fmt.Sprintf("%02x", di))

	if len(params) < 1 {
		return b
	}

	if ts, ok := params[0].(int64); ok {
		b.WriteByte('_')
		b.WriteString(time.Unix(ts, 0).Format("200601"))
	}

	return b

}

func ProvideMonthlyDbSelector() *MonthlyDbSelector {
	return &MonthlyDbSelector{}
}

type AnnuallyDbSelector struct {
}

func (se *AnnuallyDbSelector) Select(di int, dbname string, params ...interface{}) *bytes.Buffer {

	b := AcquireBuffer()
	b.WriteString(dbname)
	b.WriteByte('_')
	b.WriteString(fmt.Sprintf("%02x", di))

	if len(params) < 1 {
		return b
	}

	if ts, ok := params[0].(int64); ok {
		b.WriteByte('_')
		b.WriteString(time.Unix(ts, 0).Format("2006"))
	}

	return b

}

func ProvideAnnuallyDbSelector() *AnnuallyDbSelector {
	return &AnnuallyDbSelector{}
}
