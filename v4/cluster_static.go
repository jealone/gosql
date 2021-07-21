package gosql

type StaticCluster struct {
	Sharding Sharding
	Shards   []Shard
	Config   *Config
	Executor DbExecutor
}

func (s *StaticCluster) Write(table string, key []byte, handler DBHandler, params ...interface{}) {
	di, ti, t := s.Sharding.Select(table, key, params...)
	defer ReleaseBuffer(t)
	s.Executor(s.Shards[di].GetMaster(), handler, ti, t)
}

func (s *StaticCluster) Read(table string, key []byte, handler DBHandler, params ...interface{}) {
	di, ti, t := s.Sharding.Select(table, key, params...)
	defer ReleaseBuffer(t)
	s.Executor(s.Shards[di].GetReplica(), handler, ti, t)
}

func (s *StaticCluster) TableSelector() TableSelector {
	return s.Sharding.TableSelector()
}

func (s *StaticCluster) GetShardsTotal() int {
	return len(s.Shards)
}
