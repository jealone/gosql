package gosql

import "github.com/jealone/distributed"

func NewReplicationConfig(conf *YamlReplicationConfig) *ReplicationConfig {
	return &ReplicationConfig{
		*conf,
	}
}

type ReplicationConfig struct {
	YamlReplicationConfig
}

func (r *ReplicationConfig) GetType() string {
	return r.Type
}

func NewReplication(conf ReplicationConfiger) Replication {
	switch conf.GetType() {
	case "roundrobin":
		return distributed.NewRoundRobin()
	default:
		return &NopReplication{}
	}
}

type NopReplication struct {
}

func (lb *NopReplication) Replicate(int) int {
	return 0
}
