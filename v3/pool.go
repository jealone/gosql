package gosql

import (
	"bytes"
	"strings"
	"sync"
)

func CreatePool(f func() interface{}) func() *sync.Pool {
	var once sync.Once
	var pool *sync.Pool
	return func() *sync.Pool {

		once.Do(func() {
			pool = &sync.Pool{New: f}
		})
		return pool
	}
}

var (
	GetBufferPool = CreatePool(func() interface{} {
		var b bytes.Buffer
		return &b
	})
)

func AcquireBuffer() *bytes.Buffer {
	if b, ok := GetBufferPool().Get().(*bytes.Buffer); ok {
		return b
	}
	var b bytes.Buffer
	return &b
}

func ReleaseBuffer(b *bytes.Buffer) {
	b.Reset()
	GetBufferPool().Put(b)
}

var (
	GetStringsBuilderPool = CreatePool(func() interface{} {
		var b bytes.Buffer
		return &b
	})
)

func AcquireStringsBuilder() *strings.Builder {
	if b, ok := GetStringsBuilderPool().Get().(*strings.Builder); ok {
		return b
	}
	var b strings.Builder
	return &b
}

func ReleaseStringsBuilder(b *strings.Builder) {
	b.Reset()
	GetStringsBuilderPool().Put(b)
}
