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
	mutex    sync.Mutex
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mutex.Lock()
	listItem, ok := c.items[key]
	if ok {
		listItem.Value = value
		c.queue.MoveToFront(listItem)
	} else {
		if c.capacity == c.queue.Len() {
			c.queue.Remove(c.queue.Back())
		}
		c.items[key] = c.queue.PushFront(value)
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
	c.mutex.Unlock()
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
