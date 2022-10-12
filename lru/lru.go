package lru

import (
	"container/list"
	"errors"
	"sync"
)

type EvictCallback func(key interface{}, value interface{})

type Cache struct {
	size      int
	evictList *list.List
	//items     map[interface{}]*list.Element
	items   *sync.Map
	onEvict EvictCallback
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
		items:     &sync.Map{},
		onEvict:   onEvict,
	}
	return c, nil
}

func (c *Cache) Purge() {
	c.items.Range(func(k, v any) bool {
		if c.onEvict != nil {
			c.onEvict(k, v.(*list.Element).Value.(*entry).value)
		}
		c.items.Delete(k)
		return true
	})

	c.evictList.Init()
}

func (c *Cache) Add(key, value interface{}) (evicted bool) {
	if v, ok := c.items.Load(key); ok {
		ent := v.(*list.Element)
		c.evictList.MoveToFront(ent)
		ent.Value.(*entry).value = value
	}

	ent := &entry{key, value}
	entry_ := c.evictList.PushFront(ent)
	c.items.Store(key, entry_)

	evict := c.evictList.Len() > c.size
	if evict {
		c.removeOldest()
	}
	return evict
}

func (c *Cache) Get(key interface{}) (value interface{}, ok bool) {
	if v, ok_ := c.items.Load(key); ok_ {
		ent := v.(*list.Element)
		c.evictList.MoveToFront(ent)
		return ent.Value.(*entry).value, true
	}
	return
}

func (c *Cache) Contains(key interface{}) (ok bool) {
	_, ok = c.items.Load(key)
	return ok
}

func (c *Cache) Peek(key interface{}) (value interface{}, ok bool) {
	if v, ok_ := c.items.Load(key); ok_ {
		ent := v.(*list.Element)
		return ent.Value.(*entry).value, true
	}
	return nil, false
}

func (c *Cache) Remove(key interface{}) (present bool) {
	if v, ok_ := c.items.Load(key); ok_ {
		ent := v.(*list.Element)
		c.removeElement(ent, true)
		return true
	}
	return false
}

func (c *Cache) RemoveWithoutCallback(key interface{}) (present bool) {
	if v, ok_ := c.items.Load(key); ok_ {
		ent := v.(*list.Element)
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
	var keys []interface{}
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
	c.items.Delete(kv.key)
	if needEvictCallback && c.onEvict != nil {
		c.onEvict(kv.key, kv.value)
	}
}

func (c *Cache) ForEach(callback func(value interface{})) {
	c.items.Range(func(key, value any) bool {
		ent := value.(*list.Element)
		callback(ent.Value.(*entry).value)
		return true
	})
}
