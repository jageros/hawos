package lru

import (
	"container/list"
	"errors"
)

type EvictCallback func(key interface{}, value interface{})

type Cache struct {
	size      int
	evictList *list.List
	items     map[interface{}]*list.Element
	onEvict   EvictCallback
}

type entry struct {
	key   interface{}
	value interface{}
}

func NewCache(size int, onEvict EvictCallback) (*Cache, error) {
	if size <= 0 {
		return nil, errors.New("must provide a positive size")
	}
	c := &Cache{
		size:      size,
		evictList: list.New(),
		items:     make(map[interface{}]*list.Element),
		onEvict:   onEvict,
	}
	return c, nil
}

func (c *Cache) Purge() {
	for k, v := range c.items {
		if c.onEvict != nil {
			c.onEvict(k, v.Value.(*entry).value)
		}
		delete(c.items, k)
	}
	c.evictList.Init()
}

func (c *Cache) Add(key, value interface{}) (evicted bool) {
	if ent, ok := c.items[key]; ok {
		c.evictList.MoveToFront(ent)
		ent.Value.(*entry).value = value
		return false
	}

	ent := &entry{key, value}
	entry := c.evictList.PushFront(ent)
	c.items[key] = entry

	evict := c.evictList.Len() > c.size
	if evict {
		c.removeOldest()
	}
	return evict
}

func (c *Cache) Get(key interface{}) (value interface{}, ok bool) {
	if ent, ok := c.items[key]; ok {
		c.evictList.MoveToFront(ent)
		return ent.Value.(*entry).value, true
	}
	return
}

func (c *Cache) Contains(key interface{}) (ok bool) {
	_, ok = c.items[key]
	return ok
}

func (c *Cache) Peek(key interface{}) (value interface{}, ok bool) {
	var ent *list.Element
	if ent, ok = c.items[key]; ok {
		return ent.Value.(*entry).value, true
	}
	return nil, ok
}

func (c *Cache) Remove(key interface{}) (present bool) {
	if ent, ok := c.items[key]; ok {
		c.removeElement(ent, true)
		return true
	}
	return false
}

func (c *Cache) RemoveWithoutCallback(key interface{}) (present bool) {
	if ent, ok := c.items[key]; ok {
		c.removeElement(ent, false)
		return true
	}
	return false
}

func (c *Cache) RemoveOldest() (key interface{}, value interface{}, ok bool) {
	ent := c.evictList.Back()
	if ent != nil {
		c.removeElement(ent, true)
		kv := ent.Value.(*entry)
		return kv.key, kv.value, true
	}
	return nil, nil, false
}

func (c *Cache) GetOldest() (key interface{}, value interface{}, ok bool) {
	ent := c.evictList.Back()
	if ent != nil {
		kv := ent.Value.(*entry)
		return kv.key, kv.value, true
	}
	return nil, nil, false
}

func (c *Cache) Keys() []interface{} {
	keys := make([]interface{}, len(c.items))
	i := 0
	for ent := c.evictList.Back(); ent != nil; ent = ent.Prev() {
		keys[i] = ent.Value.(*entry).key
		i++
	}
	return keys
}

func (c *Cache) Len() int {
	return c.evictList.Len()
}

func (c *Cache) removeOldest() {
	ent := c.evictList.Back()
	if ent != nil {
		c.removeElement(ent, true)
	}
}

func (c *Cache) removeElement(e *list.Element, needEvictCallback bool) {
	c.evictList.Remove(e)
	kv := e.Value.(*entry)
	delete(c.items, kv.key)
	if needEvictCallback && c.onEvict != nil {
		c.onEvict(kv.key, kv.value)
	}
}

func (c *Cache) ForEach(callback func(value interface{})) {
	for _, ent := range c.items {
		callback(ent.Value.(*entry).value)
	}
}
