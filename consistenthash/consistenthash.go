package cache

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// Hash maps bytes to uint32
type Hash func(data []byte) uint32

// Map constains all hashed keys
type Map struct {
	hash     Hash
	replicas int            // 虚拟节点倍数
	keys     []int          // sorted 哈希环
	hashMap  map[int]string // 虚拟节点和真实节点的映射表
}

func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// add some keys to the hash
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}
	sort.Ints(m.keys)
}

func (m *Map) Get(key string) string {
	l := len(m.keys)
	if l == 0 {
		return ""
	}

	hash := int(m.hash([]byte(key)))
	// 顺时针找到第一个匹配虚拟节点的下标
	idx := sort.Search(l, func(i int) bool {
		return m.keys[i] >= hash
	})

	return m.hashMap[m.keys[idx%l]]
}
