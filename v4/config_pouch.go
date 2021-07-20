package gosql

import (
	"io"

	"github.com/google/wire"
	"gopkg.in/yaml.v3"
)

type (
	YamlDecoder = yaml.Decoder
)

var (
	defaultYamlConfigSet = wire.NewSet(
		ProvideYamlConfig,
		wire.Bind(new(Decoder), new(*YamlDecoder)),
		ProvideYamlDecoder,
	)

	DefaultConfigSet = wire.NewSet(
		defaultYamlConfigSet,
		ProvideConfig,
	)
)

func ProvideYamlDecoder(rd io.Reader) *YamlDecoder {
	return yaml.NewDecoder(rd)
}
