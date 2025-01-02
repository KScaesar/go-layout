package utility

import (
	"strings"
	"sync"
)

func NewKeyBuilder() KeyBuilder {
	return KeyBuilder{
		pool: &sync.Pool{
			New: func() interface{} {
				return &strings.Builder{}
			},
		},
	}
}

type KeyBuilder struct {
	pool *sync.Pool
	buf  *strings.Builder
}

func (k KeyBuilder) get() *strings.Builder {
	return k.pool.Get().(*strings.Builder)
}

func (k KeyBuilder) put() {
	k.buf.Reset()
	k.pool.Put(k.buf)
	k.buf = nil
}

func (k KeyBuilder) InitWithVersion(version string) KeyBuilder {
	if version == "" {
		return k.Init()
	}

	buf := k.get()
	buf.WriteString(version)
	return KeyBuilder{
		pool: k.pool,
		buf:  buf,
	}
}

func (k KeyBuilder) Init() KeyBuilder {
	buf := k.get()
	return KeyBuilder{
		pool: k.pool,
		buf:  buf,
	}
}

func (k KeyBuilder) BuildString(elements ...string) string {
	defer k.put()
	for _, element := range elements {
		if element == "" {
			element = "empty"
		}
		k.buf.WriteString(element)
	}
	return k.buf.String()
}
