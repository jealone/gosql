package gosql

import (
	"bytes"
	"fmt"
	"time"
)

type TableSelector interface {
	Select(string, []byte, ...interface{}) (int, int, *bytes.Buffer)
	Pick(string, int, ...interface{}) *bytes.Buffer
	Total(string) int
}

type StandardTableSelector struct {
	conf ShardingConfig
	m    map[string]int
	dp   Partition
}

func (s *StandardTableSelector) Select(table string, key []byte, params ...interface{}) (int, int, *bytes.Buffer) {
	b := AcquireBuffer()
	b.WriteString(table)
	total := s.Total(table)
	var index int
	if total > 1 {
		b.WriteByte('_')
		index = s.dp.Partition(key, total)
		b.WriteString(fmt.Sprintf("%02x", index))
	}

	return index, total, b
}

func (s *StandardTableSelector) Pick(table string, node int, params ...interface{}) *bytes.Buffer {
	b := AcquireBuffer()
	b.WriteString(table)
	total := s.Total(table)

	if total == 1 {
		return b
	}

	b.WriteByte('_')
	b.WriteString(fmt.Sprintf("%02x", node))
	return b
}

func (s *StandardTableSelector) Total(table string) int {
	if t, ok := s.m[table]; ok {
		return t
	} else {
		return s.conf.GetTotal()
	}
}

func ProvideStandardTableSelector(conf *Config, dp Partition) *StandardTableSelector {

	m := make(map[string]int, len(conf.GetShardingConfig().GetTableConfig()))
	for _, c := range conf.GetShardingConfig().GetTableConfig() {
		m[c.GetTableName()] = c.GetTotal()
	}

	return &StandardTableSelector{
		conf: conf.GetShardingConfig(),
		m:    m,
		dp:   dp,
	}
}

type DailyTableSelector struct {
	conf ShardingConfig
	m    map[string]int
	dp   Partition
}

func (s *DailyTableSelector) Select(table string, key []byte, params ...interface{}) (int, int, *bytes.Buffer) {

	b := AcquireBuffer()
	b.WriteString(table)
	total := s.Total(table)

	var index int
	if total > 1 {
		b.WriteByte('_')
		index = s.dp.Partition(key, total)
		b.WriteString(fmt.Sprintf("%02x", index))
	}

	if len(params) < 1 {
		return index, total, b
	}

	if ts, ok := params[0].(int64); ok {
		b.WriteByte('_')
		b.WriteString(time.Unix(ts, 0).Format("20060102"))
	}
	return index, total, b
}

func (s *DailyTableSelector) Pick(table string, node int, params ...interface{}) *bytes.Buffer {
	b := AcquireBuffer()
	b.WriteString(table)
	total := s.Total(table)
	if total > 1 {
		b.WriteByte('_')
		b.WriteString(fmt.Sprintf("%02x", node))
	}

	if ts, ok := params[0].(int64); ok {
		b.WriteByte('_')
		b.WriteString(time.Unix(ts, 0).Format("20060102"))
	}
	return b
}

func (s *DailyTableSelector) Total(table string) int {
	if t, ok := s.m[table]; ok {
		return t
	} else {
		return s.conf.GetTotal()
	}
}

func ProvideDailyTableSelector(conf *Config, dp Partition) *DailyTableSelector {

	m := make(map[string]int, len(conf.GetShardingConfig().GetTableConfig()))
	for _, c := range conf.GetShardingConfig().GetTableConfig() {
		m[c.GetTableName()] = c.GetTotal()
	}

	return &DailyTableSelector{
		conf: conf.GetShardingConfig(),
		m:    m,
		dp:   dp,
	}
}

type MonthlyTableSelector struct {
	conf ShardingConfig
	m    map[string]int
	dp   Partition
}

func (s *MonthlyTableSelector) Select(table string, key []byte, params ...interface{}) (int, int, *bytes.Buffer) {

	b := AcquireBuffer()
	b.WriteString(table)

	total := s.Total(table)
	var index int
	if total > 1 {
		b.WriteByte('_')
		index = s.dp.Partition(key, total)
		b.WriteString(fmt.Sprintf("%02x", index))
	}

	if len(params) < 1 {
		return index, total, b
	}

	if ts, ok := params[0].(int64); ok {
		b.WriteByte('_')
		b.WriteString(time.Unix(ts, 0).Format("200601"))
	}
	return index, total, b
}

func (s *MonthlyTableSelector) Pick(table string, node int, params ...interface{}) *bytes.Buffer {
	b := AcquireBuffer()
	b.WriteString(table)

	total := s.Total(table)
	if total > 1 {
		b.WriteByte('_')
		b.WriteString(fmt.Sprintf("%02x", node))
	}

	if ts, ok := params[0].(int64); ok {
		b.WriteByte('_')
		b.WriteString(time.Unix(ts, 0).Format("200601"))
	}
	return b
}

func (s *MonthlyTableSelector) Total(table string) int {
	if t, ok := s.m[table]; ok {
		return t
	} else {
		return s.conf.GetTotal()
	}
}

func ProvideMonthlyTableSelector(conf *Config, dp Partition) *MonthlyTableSelector {

	m := make(map[string]int, len(conf.GetShardingConfig().GetTableConfig()))
	for _, c := range conf.GetShardingConfig().GetTableConfig() {
		m[c.GetTableName()] = c.GetTotal()
	}

	return &MonthlyTableSelector{
		conf: conf.GetShardingConfig(),
		m:    m,
		dp:   dp,
	}
}

type AnnuallyTableSelector struct {
	conf ShardingConfig
	m    map[string]int
	dp   Partition
}

func (s *AnnuallyTableSelector) Select(table string, key []byte, params ...interface{}) (int, int, *bytes.Buffer) {

	b := AcquireBuffer()
	b.WriteString(table)

	total := s.Total(table)
	var index int
	if total > 1 {
		b.WriteByte('_')
		index = s.dp.Partition(key, total)
		b.WriteString(fmt.Sprintf("%02x", index))
	}

	if len(params) < 1 {
		return index, total, b
	}

	if ts, ok := params[0].(int64); ok {
		b.WriteByte('_')
		b.WriteString(time.Unix(ts, 0).Format("2006"))
	}
	return index, total, b
}

func (s *AnnuallyTableSelector) Pick(table string, node int, params ...interface{}) *bytes.Buffer {
	b := AcquireBuffer()
	b.WriteString(table)

	total := s.Total(table)
	if total > 1 {
		b.WriteByte('_')
		b.WriteString(fmt.Sprintf("%02x", node))
	}

	if ts, ok := params[0].(int64); ok {
		b.WriteByte('_')
		b.WriteString(time.Unix(ts, 0).Format("2006"))
	}

	return b
}

func (s *AnnuallyTableSelector) Total(table string) int {
	if t, ok := s.m[table]; ok {
		return t
	} else {
		return s.conf.GetTotal()
	}
}

func ProvideAnnuallyTableSelector(conf *Config, dp Partition) *AnnuallyTableSelector {

	m := make(map[string]int, len(conf.GetShardingConfig().GetTableConfig()))
	for _, c := range conf.GetShardingConfig().GetTableConfig() {
		m[c.GetTableName()] = c.GetTotal()
	}

	return &AnnuallyTableSelector{
		conf: conf.GetShardingConfig(),
		m:    m,
		dp:   dp,
	}
}
