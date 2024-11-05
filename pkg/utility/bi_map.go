package utility

import (
	"fmt"
	"unique"
)

func NewBiMap[K, V comparable]() *BiMap[K, V] {
	return &BiMap[K, V]{
		keyMap: make(map[unique.Handle[K]]unique.Handle[V]),
		valMap: make(map[unique.Handle[V]]unique.Handle[K]),
		keys:   make(map[K]V),
		values: make(map[V]K),
	}
}

// BiMap represents a one-to-one bidirectional map
type BiMap[K, V comparable] struct {
	keyMap map[unique.Handle[K]]unique.Handle[V]
	valMap map[unique.Handle[V]]unique.Handle[K]

	keys   map[K]V
	values map[V]K
}

func (b *BiMap[K, V]) MustSet(key K, val V) *BiMap[K, V] {
	b, err := b.Set(key, val)
	if err != nil {
		panic(err)
	}
	return b
}

func (b *BiMap[K, V]) Set(key K, val V) (*BiMap[K, V], error) {
	k, v := unique.Make(key), unique.Make(val)
	if _, ok := b.keyMap[k]; ok {
		return nil, fmt.Errorf("duplicate key: %v", key)
	}
	if _, ok := b.valMap[v]; ok {
		return nil, fmt.Errorf("duplicate val: %v", val)
	}

	b.keyMap[k] = v
	b.valMap[v] = k

	b.keys[k.Value()] = v.Value()
	b.values[v.Value()] = k.Value()

	return b, nil
}

func (b *BiMap[K, V]) GetByKey(key K) (V, bool) {
	val, ok := b.keyMap[unique.Make(key)]
	if !ok {
		var empty V
		return empty, false
	}
	return val.Value(), true
}

func (b *BiMap[K, V]) GetByValue(val V) (K, bool) {
	key, ok := b.valMap[unique.Make(val)]
	if !ok {
		var empty K
		return empty, false
	}
	return key.Value(), true
}

func (b *BiMap[K, V]) KeyMapping() map[K]V {
	return b.keys
}

func (b *BiMap[K, V]) ValueMapping() map[V]K {
	return b.values
}
