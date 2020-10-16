package gosql

import (
	"fmt"
	"sync"

	"github.com/jealone/distributed"
)

// 分库
func NewPartitionConfig(p *YamlPartitionConfig) *PartitionConfig {
	return &PartitionConfig{
		Type:    p.Type,
		DbTotal: p.DbTotal,
		Table:   NewTablePartitionConfig(p.Table, p.DbTotal),
	}
}

type PartitionConfig struct {
	Type    string
	DbTotal int
	Table   *TablePartitionConfig
}

func (p *PartitionConfig) GetTablePartitionConfig() *TablePartitionConfig {
	return p.Table
}

func (p *PartitionConfig) GetDbTotal() int {
	if p.DbTotal < 1 {
		return 1
	}
	return p.DbTotal
}

func (p *PartitionConfig) GetType() string {
	return p.Type
}

func NewPartition(conf PartitionConfiger) (DbPartitioner, error) {

	switch conf.GetType() {
	case "modular":
		p := &DbPartition{
			Total: conf.GetDbTotal(),
		}

		p.tablePartitions = make(map[string]TablePartition)

		for _, c := range conf.GetTablePartitionConfig().GetTablePartitionNodeConfig() {
			p.tablePartitions[c.GetName()] = NewTablePartition(c)
		}

		return p, nil
	default:
		return GetNopDbPartition(), NopDbPartitionError
	}
}

type DbPartitioner interface {
	Pick(prefix string, index int) string
	Partition([]byte, string) int
	TableSelector
}

type TableSelector interface {
	Select([]byte, string) string
}

type DbPartition struct {
	Total           int
	tablePartitions map[string]TablePartition
}

func (p *DbPartition) Partition(key []byte, table string) int {
	tp := p.getTablePartition(table)
	return tp.Partition(key) / tp.GetTableChunkTotal()
}

func (p *DbPartition) Select(key []byte, table string) string {
	return p.getTablePartition(table).Select(key, table)
}

func (p *DbPartition) Pick(prefix string, num int) string {
	return DividesGenerate(prefix, num)
}

func (p *DbPartition) getTablePartition(table string) TablePartition {
	if tp, ok := p.tablePartitions[table]; ok {
		return tp
	}
	return GetNopTablePartition()
}

func createNopDbPartition() func() *NopDbPartition {
	var (
		once sync.Once
		p    *NopDbPartition
	)

	return func() *NopDbPartition {
		once.Do(func() {
			p = &NopDbPartition{}
		})
		return p
	}
}

var GetNopDbPartition = createNopDbPartition()

type NopDbPartition struct {
}

func (p *NopDbPartition) Partition(bytes []byte, s string) int {
	return 0
}

func (p *NopDbPartition) Pick(prefix string, num int) string {
	return prefix
}

func (p *NopDbPartition) Select(key []byte, table string) string {
	return table
}

func (p *NopDbPartition) GetTablePartition(table string) TablePartition {
	return &NopTablePartition{}
}

// 分表
func NewTablePartitionConfig(conf []*YamlTablePartitionConfig, dbTotal int) *TablePartitionConfig {
	var config []*TablePartitionNodeConfig

	for _, c := range conf {

		config = append(config, &TablePartitionNodeConfig{
			DbTotal: dbTotal,
			Name:    c.Name,
			Total:   c.Total,
		})
	}

	return &TablePartitionConfig{
		DbTotal:     dbTotal,
		TableConfig: config,
	}
}

type TablePartitionConfig struct {
	DbTotal     int
	TableConfig []*TablePartitionNodeConfig
}

func (c *TablePartitionConfig) GetTablePartitionNodeConfig() []TablePartitionConfiger {
	var config []TablePartitionConfiger
	for _, c := range c.TableConfig {
		config = append(config, c)
	}
	return config
}

type TablePartitionNodeConfig struct {
	DbTotal int
	Name    string
	Total   int
}

func (t *TablePartitionNodeConfig) GetDbTotal() int {
	if t.DbTotal <= 1 {
		return 1
	}
	return t.DbTotal
}

func (t *TablePartitionNodeConfig) GetTableTotal() int {
	if t.Total <= 1 {
		return 1
	}
	return t.Total
}

func (t *TablePartitionNodeConfig) GetName() string {
	return t.Name
}

type TablePartition interface {
	Select(key []byte, table string) string
	Partition(key []byte) int
	GetTableTotal() int
	GetTableChunkTotal() int
}

func createNopTablePartition() func() *NopTablePartition {
	var (
		once sync.Once
		p    *NopTablePartition
	)

	return func() *NopTablePartition {
		once.Do(func() {
			p = &NopTablePartition{}
		})
		return p
	}
}

var GetNopTablePartition = createNopTablePartition()

type NopTablePartition struct {
}

func (p *NopTablePartition) Select(key []byte, table string) string {
	return table
}

func (p *NopTablePartition) Partition(key []byte) int {
	return 0
}

func (p *NopTablePartition) GetTableTotal() int {
	return 1
}

func (p *NopTablePartition) GetTableChunkTotal() int {
	return 1
}

type PreloadTablePartition struct {
	tables          []string
	Table           string
	TableTotal      int
	TableChunkTotal int
	partition       Partition
}

func (p *PreloadTablePartition) GetTableTotal() int {
	return p.TableTotal
}

func (p *PreloadTablePartition) GetTableChunkTotal() int {
	return p.TableChunkTotal
}

func (p *PreloadTablePartition) Partition(key []byte) int {
	return p.partition.Partition(key, p.GetTableTotal())
}

func (p *PreloadTablePartition) Select(key []byte, table string) string {
	if p.Table != table {
		return table
	}
	return p.pick(p.Partition(key))
}

func (p *PreloadTablePartition) pick(index int) string {
	return p.tables[index]
}

func NewTablePartition(conf TablePartitionConfiger) TablePartition {

	if conf.GetTableTotal() <= 1 {
		return GetNopTablePartition()
	}

	tp := &PreloadTablePartition{
		Table:           conf.GetName(),
		partition:       distributed.ModularHashIEEE(),
		TableTotal:      conf.GetTableTotal() * conf.GetDbTotal(),
		TableChunkTotal: conf.GetTableTotal(),
	}

	var tables []string

	for i := 0; i < tp.GetTableTotal(); i++ {
		tables = append(tables, DividesGenerate(conf.GetName(), i))
	}

	tp.tables = tables
	return tp
}

func DividesGenerate(prefix string, num int) string {
	b := AcquireBuffer()
	defer ReleaseBuffer(b)

	b.WriteString(prefix)
	b.WriteByte('_')
	b.WriteString(fmt.Sprintf("%02x", num))

	return b.String()
}
