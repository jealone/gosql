package gosql

/*
import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jealone/distributed"
)

func shardConf() *ShardConfig {
	return NewShardConfig(&YamlShardConfig{
		Master: &YamlMasterConfig{
			Host: "127.0.0.1",
			YamlNodeConfig: &YamlNodeConfig{
				User:     "root",
				Password: "admin",
			},
		},
		Replicas: &YamlReplicasConfig{
			Hosts: []string{
				"127.0.0.1",
				"127.0.0.1",
			},
			YamlNodeConfig: &YamlNodeConfig{
				User:     "root",
				Password: "admin",
			},
		},
	}, &YamlDBConfig{
		Driver: "mysql",
		Dbname: "test",
		Conn: &YamlConnConfig{
			MaxLifetime:  0,
			MaxOpenConns: 0,
			MaxIdleConns: 0,
		},
	})

}

func TestNewShardConfig(t *testing.T) {
	db := &YamlDBConfig{
		Driver: "mysql",
		Dbname: "test",
		Conn: &YamlConnConfig{
			MaxLifetime:  0,
			MaxOpenConns: 0,
			MaxIdleConns: 0,
		},
	}
	shard := &YamlShardConfig{
		Master: &YamlMasterConfig{
			Host: "127.0.0.1",
			YamlNodeConfig: &YamlNodeConfig{
				User:     "root",
				Password: "admin",
			},
		},
		Replicas: &YamlReplicasConfig{
			Hosts: []string{
				"127.0.0.1",
				"127.0.0.1",
			},
			YamlNodeConfig: &YamlNodeConfig{
				User:     "root",
				Password: "admin",
			},
		},
	}

	conf := NewShardConfig(shard, db)

	conf.SetDbname("test2")

	fmt.Println(conf.GetReplicaUrls())
	if db.Dbname == conf.GetDbname() {
		t.Fatalf("copy config fail")
	}
}

func TestShards(t *testing.T) {

	s := NewShard(shardConf(), AddExecutor(func(db *sql.DB, handler DBHandler, tp TablePartition) {
		err := handler(db, tp)
		if nil != err {
			t.Log(err)
		}
	}), AddTablePartition(NewPreloadTP(distributed.ModularHashIEEE())))

	var rows *sql.Rows
	var err error

	s.Read(func(db *sql.DB, tp TablePartition) error {
		rows, err = db.Query("SELECT * FROM `test`")
		if nil != err {
			return err
		}
		return nil
	})

	defer rows.Close()

	for rows.Next() {
		var (
			id   int
			test string
		)

		rows.Scan(&id, &test)
		fmt.Printf("%d => %s\n", id, test)
	}

}

func BenchmarkShard(b *testing.B) {

	b.ReportAllocs()

	s := NewShard(shardConf(), AddExecutor(func(db *sql.DB, handler DBHandler, tp TablePartition) {
		err := handler(db, tp)
		if nil != err {
			b.Log(err)
		}
	}), AddTablePartition(NewPreloadTP(distributed.ModularHashIEEE())))

	var (
		id   int
		test string
	)

	for i := 0; i < b.N; i++ {
		s.Read(func(db *sql.DB, partition TablePartition) error {
			row := db.QueryRow("SELECT * FROM `test` WHERE `id` = ?", 1)
			return row.Scan(&id, &test)
		})
	}
}


*/
