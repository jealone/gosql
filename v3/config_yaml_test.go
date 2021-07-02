package gosql

import (
	"os"
	"reflect"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestParseYamlConfigFile(t *testing.T) {
	type args struct {
		path string
	}

	tests := []struct {
		name string
		args args
		want *YamlConfig
	}{
		{
			name: "normal",
			args: args{
				path: "config/demo.yml",
			},
			want: wantYamlConfig(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseYamlConfigFile(tt.args.path)
			if nil != err {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(got.Cluster, tt.want.Cluster) {
				t.Errorf("got Cluster = %+v, want %+v", got.Cluster, tt.want.Cluster)
			}

			if !reflect.DeepEqual(got.Db, tt.want.Db) {
				t.Errorf("got Db = %+v, want %+v", got.Db, tt.want.Db)
			}

			if len(got.Shards) != len(tt.want.Shards) {
				t.Errorf("got Shards len = %+v, want %+v", len(got.Shards), len(tt.want.Shards))
				t.FailNow()
			}

			for i := 0; i < len(tt.want.Shards); i++ {

				if !reflect.DeepEqual(got.Shards[i].Master, tt.want.Shards[i].Master) {
					t.Errorf("node %d got shards Master = %+v, want %+v", i, got.Shards[i].Master, tt.want.Shards[i].Master)
				}

				if !reflect.DeepEqual(got.Shards[i].Replicas, tt.want.Shards[i].Replicas) {
					t.Errorf("node %d got shards Replicas = %+v, want %+v", i, got.Shards[i].Replicas, tt.want.Shards[i].Replicas)
				}
			}
		})
	}
}

func TestParseYamlConfig(t *testing.T) {

	type args struct {
		dec Decoder
	}

	file, _ := os.Open("config/demo.yml")

	tests := []struct {
		name string
		args args
		want *YamlConfig
	}{
		{
			name: "normal",
			args: args{
				dec: yaml.NewDecoder(file),
			},
			want: wantYamlConfig(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseYamlConfig(tt.args.dec)
			if nil != err {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(got.Cluster, tt.want.Cluster) {
				t.Errorf("got Cluster = %+v, want %+v", got.Cluster, tt.want.Cluster)
			}

			if !reflect.DeepEqual(got.Db, tt.want.Db) {
				t.Errorf("got Db = %+v, want %+v", got.Db, tt.want.Db)
			}

			if len(got.Shards) != len(tt.want.Shards) {
				t.Errorf("got Shards len = %+v, want %+v", len(got.Shards), len(tt.want.Shards))
				t.FailNow()
			}

			for i := 0; i < len(tt.want.Shards); i++ {

				if !reflect.DeepEqual(got.Shards[i].Master, tt.want.Shards[i].Master) {
					t.Errorf("node %d got shards Master = %+v, want %+v", i, got.Shards[i].Master, tt.want.Shards[i].Master)
				}

				if !reflect.DeepEqual(got.Shards[i].Replicas, tt.want.Shards[i].Replicas) {
					t.Errorf("node %d got shards Replicas = %+v, want %+v", i, got.Shards[i].Replicas, tt.want.Shards[i].Replicas)
				}
			}

		})
	}
}

func wantYamlConfig() *YamlConfig {
	return &YamlConfig{
		Cluster: YamlClusterConfig{
			Replication: YamlReplicationConfig{
				Type: "roundrobin",
			},
			Partition: YamlPartitionConfig{
				Type:    "modular",
				DbTotal: 2,
				Table: []*YamlTablePartitionConfig{
					{
						Name:  "test",
						Total: 8,
					},
					{
						Name:  "user",
						Total: 4,
					},
				},
			},
		},
		Db: YamlDBConfig{
			Driver: "mysql",
			Dbname: "test",
			Conn:   YamlConnConfig{},
		},
		Shards: []*YamlShardConfig{
			{
				Master: YamlMasterConfig{
					Host: "127.0.0.1:3306",
					YamlNodeConfig: YamlNodeConfig{
						User:     "root",
						Password: "admin",
						Params: map[string]string{
							"charset": "utf8",
						},
					},
				},
				Replicas: YamlReplicasConfig{
					Hosts: []string{
						"127.0.0.1:3306",
					},
					YamlNodeConfig: YamlNodeConfig{
						User:     "root",
						Password: "admin",
						Params: map[string]string{
							"charset": "utf8",
						},
					},
				},
			},
			{
				Master: YamlMasterConfig{
					Host: "127.0.0.1:3306",
					YamlNodeConfig: YamlNodeConfig{
						User:     "root",
						Password: "admin",
						Params: map[string]string{
							"charset": "utf8",
						},
					},
				},
				Replicas: YamlReplicasConfig{
					Hosts: []string{
						"127.0.0.1:3306",
					},
					YamlNodeConfig: YamlNodeConfig{
						User:     "root",
						Password: "admin",
						Params: map[string]string{
							"charset": "utf8",
						},
					},
				},
			},
		},
	}
}
