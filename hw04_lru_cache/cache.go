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
	keys     map[*ListItem]Key
	mu       *sync.Mutex
}

func (c lruCache) Set(key Key, value interface{}) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	elem := c.items[key]
	if elem == nil {
		c.queue.PushFront(value)
		c.items[key] = c.queue.Front()
		c.keys[c.queue.Front()] = key
		if c.capacity < c.queue.Len() {
			elemToRemove := c.queue.Back()
			c.queue.Remove(elemToRemove)
			keyToDelete := c.keys[elemToRemove]
			delete(c.items, keyToDelete)
			delete(c.keys, elemToRemove)
		}
		return false
	}
	elem.Value = value
	c.queue.MoveToFront(elem)
	return true
}

func (c lruCache) Get(key Key) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	elem := c.items[key]
	if elem == nil {
		return nil, false
	}
	c.queue.MoveToFront(elem)
	return elem.Value, true
}

func (c lruCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for k := range c.items {
		delete(c.items, k)
	}
	for c.queue.Len() != 0 {
		c.queue.Remove(c.queue.Back())
	}
}

func NewCache(capacity int) Cache {
	var mutex sync.Mutex
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
		keys:     make(map[*ListItem]Key),
		mu:       &mutex,
	}
}
