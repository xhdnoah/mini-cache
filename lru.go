package lru

import "container/list"

type Cache struct {
	maxBytes int64 // 允许使用最大内存
	nbytes   int64 // 当前已使用内存
	ll       *list.List
	cache    map[string]*list.Element
	// optional and executed when an entry is purged
	OnEvicted func(key string, value Value)
}

// 双向链表节点数据类型
type entry struct {
	key   string
	value Value
}

// Value use Len to count how many bytes it takes
type Value interface {
	Len() int
}

// Constructor of Cache
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

func (c *Cache) Get(key string) (value Value, ok bool) {
	if el, ok := c.cache[key]; ok {
		// 移至队尾
		c.ll.MoveToFront(el)
		kv := el.Value.(*entry)
		return kv.value, true
	}
	return
}

func (c *Cache) RemoveOldest() {
	el := c.ll.Back()
	if el != nil {
		c.ll.Remove(el)
		kv := el.Value.(*entry)
		delete(c.cache, kv.key)
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

func (c *Cache) Add(key string, value Value) {
	if el, ok := c.cache[key]; ok {
		c.ll.MoveToFront(el)
		kv := el.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		// 不存在在队尾新增节点
		el := c.ll.PushFront(&entry{key, value})
		c.cache[key] = el
		c.nbytes += int64(len(key)) + int64(value.Len())
	}
	// 更新 nbytes 如果超出则移除 LRU
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}

func (c *Cache) Len() int {
	return c.ll.Len()
}
