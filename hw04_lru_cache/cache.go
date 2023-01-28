package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	mu       *sync.RWMutex
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
		mu:       &sync.RWMutex{},
	}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, ok := c.items[key]
	if !ok {
		newItem := c.queue.PushFront(cacheItem{
			key:   key,
			value: value,
		})
		c.items[key] = newItem

		if c.queue.Len() > c.capacity {
			lastItem := c.queue.Back()
			c.queue.Remove(lastItem)
			delete(c.items, lastItem.Value.(cacheItem).key)
		}
		return false
	}

	item.Value = cacheItem{
		key:   key,
		value: value,
	}

	c.queue.MoveToFront(item)

	return true
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, ok := c.items[key]
	if !ok {
		return nil, false
	}

	c.queue.MoveToFront(item)

	return item.Value.(cacheItem).value, true
}

func (c *lruCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.queue = &list{}
	c.items = make(map[Key]*ListItem)
}
