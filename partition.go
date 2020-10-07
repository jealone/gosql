package gosql

import (
	"fmt"

	"github.com/jealone/distributed"
)

// 分库

func NewDbPartitionConfig(p *YamlPartitionConfig) *DbPartitionConfig {
	return &DbPartitionConfig{
		YamlDbPartitionConfig: p.Db,
	}
}

type DbPartitionConfig struct {
	YamlDbPartitionConfig
}

func (d *DbPartitionConfig) GetTotal() int {
	return d.Total
}

func (d *DbPartitionConfig) GetType() string {
	return d.Type
}

func NewDbPartition(conf DbPartitionConfiger) (DbPartitioner, error) {

	switch conf.GetType() {
	case "modular":
		return &DbPartition{
			Total:     conf.GetTotal(),
			Partition: distributed.ModularHashIEEE(),
		}, nil
	default:
		return &NopDbPartition{}, NopDbPartitionError
	}
}

type DbPartitioner interface {
	Select(key []byte) int
	Divide(prefix string, num int) string
}

type DbPartition struct {
	Total     int
	Partition Partition
}

func (p *DbPartition) Select(key []byte) int {
	return p.Partition.Partition(key, p.Total)
}

func (p *DbPartition) Divide(prefix string, num int) string {
	return DividesGenerate(prefix, num)
}

type NopDbPartition struct {
}

func (p *NopDbPartition) Select([]byte) int {
	return 0
}

func (p *NopDbPartition) Divide(prefix string, num int) string {
	return prefix
}

// 分表

func NewTablePartitionConfig(p *YamlPartitionConfig) *TablePartitionConfig {
	return &TablePartitionConfig{
		YamlTablePartitionConfig: p.Table,
	}
}

type TablePartitionConfig struct {
	YamlTablePartitionConfig
}

func (t *TablePartitionConfig) GetType() string {
	return t.Type
}

func (t *TablePartitionConfig) GetSpecs() []*TablePartitionSpecsConfig {
	return t.Specs
}

type (
	TablePartitionSpecsConfig = YamlTablePartitionSpecConfig
)

func (p *TablePartitionSpecsConfig) GetName() string {
	return p.Name
}

func (p *TablePartitionSpecsConfig) GetTotal() int {
	return p.Total
}

type TablePartition interface {
	Select(key []byte, prefix string) string
}

type NopTablePartition struct {
}

func (p *NopTablePartition) Select(key []byte, prefix string) string {
	return prefix
}

type PreloadTablePartition struct {
	preload   map[string][]string
	Partition Partition
}

func (p *PreloadTablePartition) Select(key []byte, prefix string) string {
	if tables, ok := p.preload[prefix]; ok {
		return tables[p.Partition.Partition(key, len(tables))]
	} else {
		return prefix
	}
}

func NewTablePartition(conf TablePartitionConfiger) TablePartition {
	switch conf.GetType() {
	case "modular":
		tp := &PreloadTablePartition{
			Partition: distributed.ModularHashIEEE(),
		}
		tp.preload = make(map[string][]string)

		for _, c := range conf.GetSpecs() {
			var t []string

			for i := 0; i < c.GetTotal(); i++ {
				t = append(t, DividesGenerate(c.GetName(), i))

			}
			tp.preload[c.GetName()] = t
		}

		return tp

	default:
		return &NopTablePartition{}
	}

}

func DividesGenerate(prefix string, num int) string {
	b := AcquireBuffer()
	defer ReleaseBuffer(b)

	b.WriteString(prefix)
	b.WriteByte('_')
	b.WriteString(fmt.Sprintf("%02x", num))

	return b.String()
}
