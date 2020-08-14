package lru

import "container/list"

type LRUCache struct {
	capacity int
	cache    map[int]*list.Element
	list     *list.List
}
type Pair struct {
	key   int
	value int
}

func NewLRUCache(capacity int) LRUCache {
	return LRUCache{
		capacity: capacity,
		list:     list.New(),
		cache:    make(map[int]*list.Element),
	}
}

func (c *LRUCache) Get(key int) int {
	if elem, ok := c.cache[key]; ok {
		c.list.MoveToFront(elem)
		return elem.Value.(Pair).value
	}
	return -1
}

func (c *LRUCache) Put(key int, value int) {
	if elem, ok := c.cache[key]; ok {
		c.list.MoveToFront(elem)
		elem.Value = Pair{key, value}
	} else {
		if c.list.Len() >= c.capacity {
			delete(c.cache, c.list.Back().Value.(Pair).key)
			c.list.Remove(c.list.Back())
		}
		c.list.PushFront(Pair{key, value})
		c.cache[key] = c.list.Front()
	}
}
