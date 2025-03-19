package util

import (
	"iter"
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

func (m *SyncMap[K, V]) Iter() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		m.m.Range(func(key, value any) bool {
			return yield(key.(K), value.(V))
		})
	}
}
