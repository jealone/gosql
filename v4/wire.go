// +build wireinject

package gosql

import (
	"io"
	"os"

	"github.com/google/wire"
)

var (
	mockYamlConfigSet = wire.NewSet(
		defaultYamlConfigSet,
		wire.Bind(new(io.Reader), new(*os.File)),
		ProvideYamlConfigFile,
	)

	mockFileConfigSet = wire.NewSet(mockYamlConfigSet, ProvideConfig)

	mockConfigSet = wire.NewSet(mockYamlConfigSet, ProvideConfig)
)

func mockDefaultYamlConfig(rd io.Reader) (*YamlConfig, error) {
	panic(wire.Build(defaultYamlConfigSet))
}

func mockFileYamlConfig(path string) (*YamlConfig, func(), error) {
	panic(wire.Build(mockYamlConfigSet))
}

func mockFileConfig(path string) (*Config, func(), error) {
	panic(wire.Build(mockConfigSet))
}

func mockDefaultConfig(rd io.Reader) (*Config, error) {
	panic(wire.Build(DefaultConfigSet))
}

func mockStaticCluster(path string) (*StaticCluster, func(), error) {
	panic(wire.Build(mockConfigSet, StaticClusterSet, StandardShardingSet, RoundRobinLBSet, DefaultDbExecutorSet, DefaultTableSelectorSet, DefaultPartitionSet))
}

func mockDailyStaticCluster(path string) (*StaticCluster, func(), error) {
	panic(wire.Build(mockConfigSet, StaticClusterSet, StandardShardingSet, RoundRobinLBSet, DefaultDbExecutorSet, DailyTableSelectorSet, DefaultPartitionSet))
}

func mockMonthlyStaticCluster(path string) (*StaticCluster, func(), error) {
	panic(wire.Build(mockConfigSet, StaticClusterSet, StandardShardingSet, RoundRobinLBSet, DefaultDbExecutorSet, MonthlyTableSelectorSet, DefaultPartitionSet))
}

func mockAnnuallyStaticCluster(path string) (*StaticCluster, func(), error) {
	panic(wire.Build(mockConfigSet, StaticClusterSet, StandardShardingSet, RoundRobinLBSet, DefaultDbExecutorSet, AnnuallyTableSelectorSet, DefaultPartitionSet))
}
