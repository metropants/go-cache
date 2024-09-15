package cache

import (
	"container/list"
	"errors"
	"fmt"
	"sync"
)

type Cache[K comparable, V any] interface {
	// Set sets an entry with the corresponding key
	Set(key K, val V) error

	// Remove removes entry from the provided key
	Remove(key K) error

	// Get returns an entry from the provided key
	Get(key K) (V, bool)

	// Exists checks whether an entry exist for the provided key
	Exists(key K) bool

	// Size returns the current size of the cache
	Size() int
}

type node[T any] struct {
	data T
	key  *list.Element
}

type MemoryCache[K comparable, V any] struct {
	mu       sync.RWMutex
	data     map[K]*node[V]
	order    *list.List
	capacity int
}

func New[K comparable, V any](capacity int) *MemoryCache[K, V] {
	return &MemoryCache[K, V]{
		data:     make(map[K]*node[V]),
		order:    list.New(),
		capacity: capacity,
	}
}

func (c *MemoryCache[K, V]) evict() error {
	el := c.order.Back()
	if el == nil {
		return errors.New("no elements to evict")
	}

	key, ok := el.Value.(K)
	if !ok {
		return errors.New("an error occurred casting K")
	}

	c.order.Remove(el)
	delete(c.data, key)
	return nil
}

func (c *MemoryCache[K, V]) Set(key K, val V) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if entry, ok := c.data[key]; !ok {
		if c.order.Len() >= c.capacity {
			err := c.evict()
			if err != nil {
				return err
			}
		}

		c.data[key] = &node[V]{
			data: val,
			key:  c.order.PushFront(key),
		}
	} else {
		c.order.MoveToFront(entry.key)
		entry.data = val
		c.data[key] = entry
	}
	return nil
}

func (c *MemoryCache[K, V]) Remove(key K) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, ok := c.data[key]
	if !ok {
		return fmt.Errorf("no entry for key: %v found", key)
	}

	delete(c.data, key)
	c.order.Remove(entry.key)
	return nil
}

func (c *MemoryCache[K, V]) Get(key K) (V, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, ok := c.data[key]
	if !ok {
		return *new(V), false
	}

	c.order.MoveToFront(entry.key)
	return entry.data, true
}

func (c *MemoryCache[K, V]) Exists(key K) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	_, ok := c.data[key]
	return ok
}

func (c *MemoryCache[K, V]) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.data)
}
