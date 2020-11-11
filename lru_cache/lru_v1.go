package lru_cache

import (
	"container/list"
	"sync"
)

type elem struct {
	data   []byte
	lruPos *list.Element
}

type lruCache struct {
	maxSize   int
	elemCount int
	data      map[interface{}]*elem //hash表
	lck       *sync.Mutex
	lru       *list.List //节点链表
}

//lru init
func NewLruCache(maxLength int) *lruCache {
	if maxLength <= 0 {
		return nil
	}
	return &lruCache{
		maxSize:   maxLength,
		elemCount: 0,
		data:      make(map[interface{}]*elem),
		lck:       new(sync.Mutex),
		lru:       list.New(),
	}
}
