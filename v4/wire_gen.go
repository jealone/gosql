// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package gosql

import (
	"github.com/google/wire"
	"github.com/jealone/distributed"
	"io"
	"os"
)

import (
	_ "github.com/go-sql-driver/mysql"
)

// Injectors from wire.go:

func mockDefaultYamlConfig(rd io.Reader) (*YamlConfig, error) {
	decoder := ProvideYamlDecoder(rd)
	yamlConfig, err := ProvideYamlConfig(decoder)
	if err != nil {
		return nil, err
	}
	return yamlConfig, nil
}

func mockFileYamlConfig(path string) (*YamlConfig, func(), error) {
	file, cleanup, err := ProvideYamlConfigFile(path)
	if err != nil {
		return nil, nil, err
	}
	decoder := ProvideYamlDecoder(file)
	yamlConfig, err := ProvideYamlConfig(decoder)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	return yamlConfig, func() {
		cleanup()
	}, nil
}

func mockFileConfig(path string) (*Config, func(), error) {
	file, cleanup, err := ProvideYamlConfigFile(path)
	if err != nil {
		return nil, nil, err
	}
	decoder := ProvideYamlDecoder(file)
	yamlConfig, err := ProvideYamlConfig(decoder)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	config, err := ProvideConfig(yamlConfig)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	return config, func() {
		cleanup()
	}, nil
}

func mockDefaultConfig(rd io.Reader) (*Config, error) {
	decoder := ProvideYamlDecoder(rd)
	yamlConfig, err := ProvideYamlConfig(decoder)
	if err != nil {
		return nil, err
	}
	config, err := ProvideConfig(yamlConfig)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func mockStaticCluster(path string) (*StaticCluster, func(), error) {
	file, cleanup, err := ProvideYamlConfigFile(path)
	if err != nil {
		return nil, nil, err
	}
	decoder := ProvideYamlDecoder(file)
	yamlConfig, err := ProvideYamlConfig(decoder)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	config, err := ProvideConfig(yamlConfig)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	modularHash := distributed.ModularHashIEEE()
	standardTableSelector := ProvideStandardTableSelector(config, modularHash)
	standardSharding := ProvideStandardSharding(config, standardTableSelector)
	roundRobin := ProvideRoundRobin()
	v := ProvideStaticShards(config, standardSharding, roundRobin)
	dbExecutor := _wireDbExecutorValue
	staticCluster := &StaticCluster{
		sharding: standardSharding,
		shards:   v,
		config:   config,
		executor: dbExecutor,
	}
	return staticCluster, func() {
		cleanup()
	}, nil
}

var (
	_wireDbExecutorValue = DefaultDBExecutor
)

func mockDailyStaticCluster(path string) (*StaticCluster, func(), error) {
	file, cleanup, err := ProvideYamlConfigFile(path)
	if err != nil {
		return nil, nil, err
	}
	decoder := ProvideYamlDecoder(file)
	yamlConfig, err := ProvideYamlConfig(decoder)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	config, err := ProvideConfig(yamlConfig)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	modularHash := distributed.ModularHashIEEE()
	dailyTableSelector := ProvideDailyTableSelector(config, modularHash)
	standardSharding := ProvideStandardSharding(config, dailyTableSelector)
	roundRobin := ProvideRoundRobin()
	v := ProvideStaticShards(config, standardSharding, roundRobin)
	dbExecutor := _wireDbExecutorValue
	staticCluster := &StaticCluster{
		sharding: standardSharding,
		shards:   v,
		config:   config,
		executor: dbExecutor,
	}
	return staticCluster, func() {
		cleanup()
	}, nil
}

func mockMonthlyStaticCluster(path string) (*StaticCluster, func(), error) {
	file, cleanup, err := ProvideYamlConfigFile(path)
	if err != nil {
		return nil, nil, err
	}
	decoder := ProvideYamlDecoder(file)
	yamlConfig, err := ProvideYamlConfig(decoder)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	config, err := ProvideConfig(yamlConfig)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	modularHash := distributed.ModularHashIEEE()
	monthlyTableSelector := ProvideMonthlyTableSelector(config, modularHash)
	standardSharding := ProvideStandardSharding(config, monthlyTableSelector)
	roundRobin := ProvideRoundRobin()
	v := ProvideStaticShards(config, standardSharding, roundRobin)
	dbExecutor := _wireDbExecutorValue
	staticCluster := &StaticCluster{
		sharding: standardSharding,
		shards:   v,
		config:   config,
		executor: dbExecutor,
	}
	return staticCluster, func() {
		cleanup()
	}, nil
}

func mockAnnuallyStaticCluster(path string) (*StaticCluster, func(), error) {
	file, cleanup, err := ProvideYamlConfigFile(path)
	if err != nil {
		return nil, nil, err
	}
	decoder := ProvideYamlDecoder(file)
	yamlConfig, err := ProvideYamlConfig(decoder)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	config, err := ProvideConfig(yamlConfig)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	modularHash := distributed.ModularHashIEEE()
	annuallyTableSelector := ProvideAnnuallyTableSelector(config, modularHash)
	standardSharding := ProvideStandardSharding(config, annuallyTableSelector)
	roundRobin := ProvideRoundRobin()
	v := ProvideStaticShards(config, standardSharding, roundRobin)
	dbExecutor := _wireDbExecutorValue
	staticCluster := &StaticCluster{
		sharding: standardSharding,
		shards:   v,
		config:   config,
		executor: dbExecutor,
	}
	return staticCluster, func() {
		cleanup()
	}, nil
}

// wire.go:

var (
	mockYamlConfigSet = wire.NewSet(
		defaultYamlConfigSet, wire.Bind(new(io.Reader), new(*os.File)), ProvideYamlConfigFile,
	)

	mockFileConfigSet = wire.NewSet(mockYamlConfigSet, ProvideConfig)

	mockConfigSet = wire.NewSet(mockYamlConfigSet, ProvideConfig)
)
