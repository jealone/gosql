package gosql

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type YamlConfig struct {
	Cluster YamlClusterConfig  `yaml:"cluster"`
	Db      YamlDBConfig       `yaml:"db"`
	Shards  []*YamlShardConfig `yaml:"shards"`
}

type YamlClusterConfig struct {
	Replication YamlReplicationConfig `yaml:"replication"`
	Partition   YamlPartitionConfig   `yaml:"partition"`
}

type YamlReplicationConfig struct {
	Type string `yaml:"type"`
}

type YamlPartitionConfig struct {
	Type    string                      `yaml:"type"`
	DbTotal int                         `yaml:"db_total"`
	Table   []*YamlTablePartitionConfig `yaml:"table"`
}

type YamlTablePartitionConfig struct {
	Name  string `yaml:"name"`
	Total int    `yaml:"total"`
}

type YamlDBConfig struct {
	Driver string         `yaml:"driver"`
	Dbname string         `yaml:"dbname"`
	Conn   YamlConnConfig `yaml:"conn"`
}

type YamlShardConfig struct {
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

type YamlConnConfig struct {
	MaxLifetime  int `yaml:"max_life_time"`
	MaxIdleConns int `yaml:"max_idle_conns"`
	MaxOpenConns int `yaml:"max_open_conns"`
}

type Decoder interface {
	Decode(interface{}) error
}

func ParseYamlConfigFile(path string) (*YamlConfig, error) {
	abs, err := filepath.Abs(path)

	if nil != err {
		return nil, err
	}

	file, err := os.Open(abs)

	if nil != err {
		return nil, err
	}

	defer file.Close()

	dec := yaml.NewDecoder(file)
	conf := &YamlConfig{}

	err = dec.Decode(conf)
	if nil != err {
		return nil, err
	}

	return conf, nil

}

func ParseYamlConfig(d Decoder) (*YamlConfig, error) {
	conf := &YamlConfig{}
	err := d.Decode(conf)
	if nil != err {
		return nil, err
	}
	return conf, nil
}
