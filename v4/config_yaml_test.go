package gosql

import (
	"io"
	"os"
	"reflect"
	"testing"
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
				path: "Config/demo.yml",
			},
			want: wantYamlConfig(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, c, err := mockFileYamlConfig(tt.args.path)
			if nil != err {
				t.Fatal(err)
			}
			defer c()

			if !reflect.DeepEqual(got.Cluster, tt.want.Cluster) {
				t.Errorf("got Cluster = %+v, want %+v", got.Cluster, tt.want.Cluster)
			}

			if !reflect.DeepEqual(got.Cluster, tt.want.Cluster) {
				t.Errorf("got Db = %+v, want %+v", got.Cluster, tt.want.Cluster)
			}

			if len(got.Shards) != len(tt.want.Shards) {
				t.Errorf("got Shards len = %+v, want %+v", len(got.Shards), len(tt.want.Shards))
				t.FailNow()
			}

			for i := 0; i < len(tt.want.Shards); i++ {

				if !reflect.DeepEqual(got.Shards[i].Master, tt.want.Shards[i].Master) {
					t.Errorf("node %d got Shards Master = %+v, want %+v", i, got.Shards[i].Master, tt.want.Shards[i].Master)
				}

				if !reflect.DeepEqual(got.Shards[i].Replicas, tt.want.Shards[i].Replicas) {
					t.Errorf("node %d got Shards Replicas = %+v, want %+v", i, got.Shards[i].Replicas, tt.want.Shards[i].Replicas)
				}
			}
		})
	}
}

func TestParseYamlConfig(t *testing.T) {

	type args struct {
		rd io.Reader
	}

	file, _ := os.Open("Config/demo.yml")

	tests := []struct {
		name string
		args args
		want *YamlConfig
	}{
		{
			name: "normal",
			args: args{
				rd: file,
			},
			want: wantYamlConfig(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mockDefaultYamlConfig(tt.args.rd)
			if nil != err {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(got.Cluster, tt.want.Cluster) {
				t.Errorf("got Cluster = %+v, want %+v", got.Cluster, tt.want.Cluster)
			}

			if !reflect.DeepEqual(got.Cluster, tt.want.Cluster) {
				t.Errorf("got Db = %+v, want %+v", got.Cluster, tt.want.Cluster)
			}

			if len(got.Shards) != len(tt.want.Shards) {
				t.Errorf("got Shards len = %+v, want %+v", len(got.Shards), len(tt.want.Shards))
				t.FailNow()
			}

			for i := 0; i < len(tt.want.Shards); i++ {

				if !reflect.DeepEqual(got.Shards[i].Master, tt.want.Shards[i].Master) {
					t.Errorf("node %d got Shards Master = %+v, want %+v", i, got.Shards[i].Master, tt.want.Shards[i].Master)
				}

				if !reflect.DeepEqual(got.Shards[i].Replicas, tt.want.Shards[i].Replicas) {
					t.Errorf("node %d got Shards Replicas = %+v, want %+v", i, got.Shards[i].Replicas, tt.want.Shards[i].Replicas)
				}
			}

		})
	}
}

func wantYamlConfig() *YamlConfig {
	return &YamlConfig{
		Cluster: YamlClusterConfig{
			YamlDBConfig: YamlDBConfig{
				Driver: "mysql",
				Dbname: "test",
				Conn:   YamlConnConfig{},
			},
			Sharding: YamlShardingConfig{
				DbTotal: 2,
				Table: []*YamlTableShardingConfig{
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
		Shards: []*YamlShardsConfig{
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
