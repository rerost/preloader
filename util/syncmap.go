package util

import (
	"sync"
)

type SyncMap[K any, V any] struct {
	m sync.Map
}

func (m *SyncMap[K, V]) Load(key K) (value *V, ok bool) {
	v, ok := m.m.Load(key)
	if !ok {
		return nil, false
	}
	res := v.(V)
	return &res, true
}

func (m *SyncMap[K, V]) Store(key K, value V) {
	m.m.Store(key, value)
}

// Range calls f sequentially for each key and value in the map
func (m *SyncMap[K, V]) Range(f func(key K, value V) bool) {
	m.m.Range(func(key, value any) bool {
		return f(key.(K), value.(V))
	})
}
