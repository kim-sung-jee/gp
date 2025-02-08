package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

type ConsistentHash struct {
	mu           sync.RWMutex
	replicas     int
	keys         []int
	hashMap      map[int]string
	hashFunction func(data []byte) uint32
}

func NewConsistentHash(replicas int, fn func(data []byte) uint32) *ConsistentHash {
	if fn == nil {
		fn = crc32.ChecksumIEEE
	}
	return &ConsistentHash{
		replicas:     replicas,
		hashMap:      make(map[int]string),
		hashFunction: fn,
	}
}

func (c *ConsistentHash) Add(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for i := 0; i < c.replicas; i++ {
		hash := int(c.hashFunction([]byte(strconv.Itoa(i) + key)))
		c.keys = append(c.keys, hash)
		c.hashMap[hash] = key
	}

	sort.Ints(c.keys)
}

func (c *ConsistentHash) Remove(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for i := 0; i < c.replicas; i++ {
		hash := int(c.hashFunction([]byte(strconv.Itoa(i) + key)))
		index := sort.SearchInts(c.keys, hash)
		if index < len(c.keys) && c.keys[index] == hash {
			c.keys = append(c.keys[:index], c.keys[index+1:]...)
		}
		delete(c.hashMap, hash)
	}
}

func (c *ConsistentHash) Get(key string) string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.keys) == 0 {
		return ""
	}

	hash := int(c.hashFunction([]byte(key)))
	index := sort.Search(len(c.keys), func(i int) bool { return c.keys[i] >= hash })

	if index == len(c.keys) {
		index = 0
	}

	return c.hashMap[c.keys[index]]
}
