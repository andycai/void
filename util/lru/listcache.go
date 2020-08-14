package lru

type Node struct {
	Key   int
	Value int
	pre   *Node
	next  *Node
}

type LRUListCache struct {
	limit   int
	HashMap map[int]*Node
	head    *Node
	end     *Node
}

func NewLRUListCache(capacity int) LRUListCache {
	LRUListCache := LRUListCache{limit: capacity}
	LRUListCache.HashMap = make(map[int]*Node, capacity)
	return LRUListCache
}

func (c *LRUListCache) Get(key int) int {
	if v, ok := c.HashMap[key]; ok {
		c.refreshNode(v)
		return v.Value
	} else {
		return -1
	}
}

func (c *LRUListCache) Put(key int, value int) {
	if v, ok := c.HashMap[key]; !ok {
		if len(c.HashMap) >= c.limit {
			oldKey := c.removeNode(c.head)
			delete(c.HashMap, oldKey)
		}
		node := Node{Key: key, Value: value}
		c.addNode(&node)
		c.HashMap[key] = &node
	} else {
		v.Value = value
		c.refreshNode(v)
	}
}

func (c *LRUListCache) refreshNode(node *Node) {
	if node == c.end {
		return
	}
	c.removeNode(node)
	c.addNode(node)
}

func (c *LRUListCache) removeNode(node *Node) int {
	if node == c.end {
		c.end = c.end.pre
		c.end.next = nil
	} else if node == c.head {
		c.head = c.head.next
		c.head.pre = nil
	} else {
		node.pre.next = node.next
		node.next.pre = node.pre
	}
	return node.Key
}

func (c *LRUListCache) addNode(node *Node) {
	if c.end != nil {
		c.end.next = node
		node.pre = c.end
		node.next = nil
	}
	c.end = node
	if c.head == nil {
		c.head = node
	}
}
