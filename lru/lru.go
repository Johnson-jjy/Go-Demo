package lru

import  "container/list"

//1.核心数据结构
//Cache is a LRU cache. It is not safe for concurrent access.
type Cache struct {
	maxBytes  int64			//max memory for use
	usedBytes int64 		//used memory
	nodeList  *list.List 	//
	cache     map[string]*list.Element
	callback  func(key string, value Value) // optional and executed when an entry is purged.
}

type entry struct {
	key string
	value Value
}

//Value use len to count how many bytes it takes
type Value interface {
	Len() int
}

//New is the Constructer of Cache
func New(maxBytes int64, callback func(string, Value)) *Cache {
	return &Cache{
		maxBytes: maxBytes,
		nodeList: list.New(),
		cache:    make(map[string]*list.Element),
		callback: callback,
	}
}

//2.1 查找功能
//Get look ups a key's value
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.nodeList.MoveToFront(ele)		//just agreed that front is the end of the team
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

//2.2 删除
//RemoveOldest removes the oldest item
func (c *Cache) RemoveOldest() {
	ele := c.nodeList.Back()
	if ele != nil {
		c.nodeList.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.usedBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.callback != nil {
			c.callback(kv.key, kv.value)
		}
	}
}

//2.3 新增/修改
//Add adds a value to the cache.
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		c.nodeList.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.usedBytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		ele := c.nodeList.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.usedBytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.maxBytes < c.usedBytes {
		c.RemoveOldest()
	}
}

//Len the number of cache entries
func (c *Cache) Len() int {
	return c.nodeList.Len()
}