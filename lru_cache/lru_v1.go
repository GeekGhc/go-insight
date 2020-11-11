package lru_cache

import (
	"container/list"
)

type elem struct {
	lruPos *list.Element //list element
}

type LruCache struct {
	capacity int
	length   int
	dataMap  map[interface{}]*elem //hash表
	lru      *list.List            //节点链表
}

//lru init
func NewLruCache(maxLength int) *LruCache {
	if maxLength <= 0 {
		return nil
	}
	return &LruCache{
		capacity: maxLength,
		length:   0,
		dataMap:  make(map[interface{}]*elem),
		lru:      list.New(),
	}
}

func (c *LruCache) Set(key, value interface{}) {
	if v, ok := c.dataMap[key]; ok {
		//删除原先list的位置
		c.deleteLruItem(v.lruPos)
		//更新到list头部
		v.lruPos = c.syncLruItem(key)
		c.dataMap[key] = v
	}
}

func (c *LruCache) deleteLruItem(pos *list.Element) {
	//删除list
	c.lru.Remove(pos)
	//删除hashMap
	delete(c.dataMap, pos.Value)
	if c.length >= 0 {
		c.length--
	}
	return
}

func (c *LruCache) syncLruItem(key interface{}) (item *list.Element) {
	item = c.lru.PushFront(key)
	c.length++
	return
}
