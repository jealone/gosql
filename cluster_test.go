package gosql

import (
	"database/sql"
	"fmt"

	//"fmt"
	"testing"
)

func TestNewCluster(t *testing.T) {

	type args struct {
		conf *Config
	}

	c, err := ParseYamlConfigFile("config/demo.yml")

	if nil != err {
		t.Fatalf("config error %s", err)
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "normal",
			args: args{
				conf: NewConfig(c),
			},
		},
	}

	for _, tt := range tests {
		var (
			id   int
			test string
		)
		t.Run(tt.name, func(t *testing.T) {

			cluster, err := NewCluster(tt.args.conf, func(db *sql.DB, handler DBHandler, partition TablePartition) {
				err := handler(db, partition)
				if nil != err {
					t.Errorf("query error %s\n", err)
				}

			})
			if nil != err {
				t.Fatalf("new cluster error : %s", err)
			}

			cluster.Read(func(db *sql.DB, partition TablePartition) error {
				table := partition.Select([]byte("test"), "test")

				row := db.QueryRow(fmt.Sprintf("SELECT * FROM %s", table))
				return row.Scan(&id, &test)

			})
			fmt.Println(id, test)
		})
	}

}
