package gosql

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jealone/distributed"
)

type (
	Replication = distributed.Replication
	Partition   = distributed.Partition
)
