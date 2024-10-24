package utility

import (
	"fmt"
)

func NewBiMap[K, V comparable]() *BiMap[K, V] {
	return &BiMap[K, V]{
		keyMap: make(map[K]V),
		valMap: make(map[V]K),
	}
}

// BiMap represents a one-to-one bidirectional map
type BiMap[K, V comparable] struct {
	keyMap map[K]V
	valMap map[V]K
}

func (b *BiMap[K, V]) MustSet(key K, val V) *BiMap[K, V] {
	b, err := b.Set(key, val)
	if err != nil {
		panic(err)
	}
	return b
}

func (b *BiMap[K, V]) Set(key K, val V) (*BiMap[K, V], error) {
	if _, ok := b.keyMap[key]; ok {
		return nil, fmt.Errorf("duplicate key: %v", key)
	}
	if _, ok := b.valMap[val]; ok {
		return nil, fmt.Errorf("duplicate val: %v", val)
	}

	b.keyMap[key] = val
	b.valMap[val] = key

	return b, nil
}

func (b *BiMap[K, V]) GetByKey(key K) (V, bool) {
	val, ok := b.keyMap[key]
	return val, ok
}

func (b *BiMap[K, V]) GetByValue(value V) K {
	return b.valMap[value]
}

func (b *BiMap[K, V]) KeyMapping() map[K]V {
	return b.keyMap
}
