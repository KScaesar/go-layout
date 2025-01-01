package utility

import (
	"strings"
	"sync"
)

func NewKeyBuilder(prefix string) *KeyBuilder {
	return &KeyBuilder{
		pool: sync.Pool{
			New: func() interface{} {
				return &strings.Builder{}
			},
		},
		prefix: prefix,
	}
}

type KeyBuilder struct {
	pool   sync.Pool
	prefix string
}

func (k *KeyBuilder) get() *strings.Builder {
	buf := k.pool.Get().(*strings.Builder)
	if k.prefix != "" {
		buf.WriteString(k.prefix)
	}
	return buf
}

func (k *KeyBuilder) put(buf *strings.Builder) {
	buf.Reset()
	k.pool.Put(buf)
}

func (k *KeyBuilder) BuildString(elements ...string) string {
	buf := k.get()
	defer k.put(buf)
	for _, element := range elements {
		if element == "" {
			element = "empty"
		}
		buf.WriteString(element)
	}
	return buf.String()
}
