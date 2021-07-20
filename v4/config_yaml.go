package gosql

import (
	"os"
	"path/filepath"
)

type Decoder interface {
	Decode(interface{}) error
}

func ProvideYamlConfigFile(path string) (*os.File, func(), error) {
	abs, err := filepath.Abs(path)
	if nil != err {
		return nil, nil, err
	}
	file, err := os.Open(abs)
	if nil != err {
		return nil, nil, err
	}

	cleanup := func() {
		file.Close()
	}

	return file, cleanup, nil
}

func ProvideYamlConfig(dec Decoder) (*YamlConfig, error) {
	conf := &YamlConfig{}
	err := dec.Decode(conf)
	if nil != err {
		return nil, err
	}
	return conf, nil
}

type YamlConfig struct {
	Cluster YamlClusterConfig   `yaml:"cluster"`
	Shards  []*YamlShardsConfig `yaml:"shards"`
}

type YamlClusterConfig struct {
	YamlDBConfig `yaml:",inline"`
	Sharding     YamlShardingConfig `yaml:"sharding"`
}

type YamlDBConfig struct {
	Driver string         `yaml:"driver"`
	Dbname string         `yaml:"dbname"`
	Type   string         `yaml:"type"`
	Conn   YamlConnConfig `yaml:"conn"`
}

type YamlConnConfig struct {
	MaxLifetime  int `yaml:"max_life_time"`
	MaxIdleConns int `yaml:"max_idle_conns"`
	MaxOpenConns int `yaml:"max_open_conns"`
}

type YamlShardingConfig struct {
	DbTotal int                        `yaml:"total"`
	Table   []*YamlTableShardingConfig `yaml:"table"`
}

type YamlTableShardingConfig struct {
	Name  string `yaml:"name"`
	Total int    `yaml:"total"`
}

type YamlShardsConfig struct {
	Master   YamlMasterConfig   `yaml:"master"`
	Replicas YamlReplicasConfig `yaml:"replicas"`
}

type YamlMasterConfig struct {
	Host           string `yaml:"host"`
	YamlNodeConfig `yaml:",inline"`
}

type YamlReplicasConfig struct {
	Hosts          []string `yaml:"hosts"`
	YamlNodeConfig `yaml:",inline"`
}

type YamlNodeConfig struct {
	User     string            `yaml:"user"`
	Password string            `yaml:"password"`
	Protocol string            `yaml:"protocol"`
	Params   map[string]string `yaml:"params"`
}
