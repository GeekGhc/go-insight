package lru_cache

import (
	"container/list"
	"sync"
)

type value struct {
	data   []byte
	lruPos *list.Element
}

type lruCache struct {
	maxSize int
	data    map[interface{}]*value
	lck     *sync.Mutex
	lru     *list.List
}

//lru init
func NewLruCache(maxLength int) *lruCache {
	if maxLength <= 0 {
		return nil
	}
	return &lruCache{
		maxSize: maxLength,
		data:    make(map[interface{}]*value),
		lck:     new(sync.Mutex),
		lru:     list.New(),
	}
}
