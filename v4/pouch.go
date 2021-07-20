package gosql

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/wire"
	"github.com/jealone/distributed"
)

var (
	StaticClusterSet = wire.NewSet(
		wire.Bind(new(Cluster), new(*StaticCluster)),
		wire.Struct(new(StaticCluster), "*"),
		ProvideStaticShards,
	)

	StandardShardingSet = wire.NewSet(
		wire.Bind(new(Sharding), new(*StandardSharding)),
		ProvideStandardSharding,
	)

	RoundRobinLBSet = wire.NewSet(
		wire.Bind(new(Replication), new(*RoundRobinSelector)),
		ProvideRoundRobinSelector,
	)

	DefaultDbExecutorSet = wire.NewSet(
		wire.Value(DefaultDBExecutor),
	)

	DefaultTableSelectorSet = wire.NewSet(
		wire.Bind(new(TableSelector), new(*StandardTableSelector)),
		ProvideStandardTableSelector,
	)

	DailyTableSelectorSet = wire.NewSet(
		wire.Bind(new(TableSelector), new(*DailyTableSelector)),
		ProvideDailyTableSelector,
	)

	MonthlyTableSelectorSet = wire.NewSet(
		wire.Bind(new(TableSelector), new(*MonthlyTableSelector)),
		ProvideMonthlyTableSelector,
	)

	AnnuallyTableSelectorSet = wire.NewSet(
		wire.Bind(new(TableSelector), new(*AnnuallyTableSelector)),
		ProvideAnnuallyTableSelector,
	)

	DefaultPartitionSet = wire.NewSet(
		wire.Bind(new(Partition), new(*distributed.ModularHash)),
		distributed.ModularHashIEEE,
	)
)

var (
	DynamicClusterSet = wire.NewSet(
		wire.Bind(new(Cluster), new(*DynamicCluster)),
		wire.Struct(new(DynamicCluster), "*"),
		ProvideDynamicRe,
	)

	DailyDbSelectorSet = wire.NewSet(
		wire.Bind(new(DbSelector), new(*DailyDbSelector)),
		ProvideDailyDbSelector,
	)

	MonthlyDbSelectorSet = wire.NewSet(
		wire.Bind(new(DbSelector), new(*MonthlyDbSelector)),
		ProvideMonthlyDbSelector,
	)

	AnnuallyDbSelectorSet = wire.NewSet(
		wire.Bind(new(DbSelector), new(*AnnuallyDbSelector)),
		ProvideAnnuallyDbSelector,
	)
)

type Partition interface {
	Partition([]byte, int) int
}

type Replication interface {
	Replicate(int) int
}

type (
	RoundRobinSelector = distributed.RoundRobinSelector
)

func ProvideRoundRobinSelector() *RoundRobinSelector {
	return distributed.NewRoundRobin()
}
