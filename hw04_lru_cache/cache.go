package hw04lrucache

import (
	"sync"
)

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
	mutex    sync.Mutex
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mutex.Lock()
	listItem, ok := c.items[key]
	if ok {
		c.queue.MoveToFront(listItem)
		c.queue.Front().Value = value
	} else {
		if c.capacity == c.queue.Len() {
			delete(c.items, c.keys[c.queue.Back()])
			delete(c.keys, c.queue.Back())
			c.queue.Remove(c.queue.Back())
		}
		c.items[key] = c.queue.PushFront(value)
		c.keys[c.queue.Front()] = key
	}
	c.mutex.Unlock()
	return ok
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mutex.Lock()
	listItem, ok := c.items[key]
	if ok {
		c.mutex.Unlock()
		return listItem.Value, ok
	}
	c.mutex.Unlock()
	return nil, false
}

func (c *lruCache) Clear() {
	c.mutex.Lock()
	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
	c.keys = make(map[*ListItem]Key, c.capacity)
	c.mutex.Unlock()
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
		keys:     make(map[*ListItem]Key, capacity),
	}
}
