package omap

import (
	"github.com/azzz/psychic-octo-pancake/internal/list"
	"sync"
)

type Pair[K, V any] struct {
	Key   K
	Value V
}

// OMap is like a Go ma[K]V map but ordered and safe for concurrency use.
// Ordered map means the values can be indexed by key, by index in the order they were set,
// or iterated in the order.
type OMap[K comparable, V any] struct {
	*sync.RWMutex

	keys  list.List[K]
	store map[K]V
}

func New[K comparable, V any]() *OMap[K, V] {
	return &OMap[K, V]{
		RWMutex: &sync.RWMutex{},
		keys:    list.New[K](),
		store:   make(map[K]V),
	}
}

// Set sets or replaces a value under a key.
// If value existed, Set returns true.
func (m *OMap[K, V]) Set(key K, value V) bool {
	m.Lock()
	defer m.Unlock()

	_, ok := m.store[key]
	m.store[key] = value
	if !ok {
		m.keys.Push(key)
	}

	return ok
}

// Get returns a value under a key.
// If a value exists, Get returns the value and true.
// Otherwise, it returns a zero value and false.
func (m *OMap[K, V]) Get(key K) (V, bool) {
	m.RLock()
	defer m.RUnlock()

	v, ok := m.store[key]
	return v, ok
}

// Nth returns an Nth value in order it was inserted into the map.
// If a value exists, Get returns the value and true.
// Otherwise, it returns a zero value and false.
func (m *OMap[K, V]) Nth(n int) (V, bool) {
	m.RLock()
	defer m.RUnlock()

	key := m.keys.Nth(n)
	if key == nil {
		var v V
		return v, false
	}

	v, ok := m.store[key.Value()]
	return v, ok
}

// Keys return ordered keys.
func (m *OMap[K, V]) Keys() []K {
	m.RLock()
	defer m.RUnlock()

	return m.keys.Values()
}

// Pairs return ordered key-value pairs.
func (m *OMap[K, V]) Pairs() []Pair[K, V] {
	m.RLock()
	defer m.RUnlock()

	var (
		pairs = make([]Pair[K, V], 0, m.keys.Len())
		keys  = m.keys.Values()
	)

	for _, key := range keys {
		val, ok := m.store[key]
		if !ok {
			continue
		}

		pairs = append(pairs, Pair[K, V]{key, val})
	}

	return pairs
}

// Delete a value under a key.
// If value existed, Delete returns true.
func (m *OMap[K, V]) Delete(key K) bool {
	m.Lock()
	defer m.Unlock()

	_, ok := m.store[key]
	if !ok {
		return false
	}

	iter := m.keys.Iter()
	for iter.Next() {
		v, _ := iter.Value()
		if v == key {
			m.keys.Delete(iter.Idx())
			delete(m.store, key)
			break
		}
	}

	return true
}
