package lru

import "container/list"

type Cache struct {
	capacity int
	cache    map[Key]*list.Element
	l        *list.List
}

type Key interface{}
type Value interface{}
type entry struct {
	key   Key
	value Value
}

func New(capacity int) *Cache {
	return &Cache{
		capacity: capacity,
		cache:    make(map[Key]*list.Element),
		l:        list.New(),
	}
}
