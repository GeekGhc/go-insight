package v1

import (
	"container/list"
	"errors"
)

type CacheNode struct {
	key, value interface{}
}

type LruCache struct {
	capacity int
	length   int
	dataMap  map[interface{}]*list.Element //hash表
	dataList *list.List                    //节点链表
}

//lru init
func NewLruCacheV1(maxLength int) *LruCache {
	if maxLength <= 0 {
		return nil
	}
	return &LruCache{
		capacity: maxLength,
		length:   0,
		dataMap:  make(map[interface{}]*list.Element),
		dataList: list.New(),
	}
}

func (lru *LruCache) Size() int {
	return lru.dataList.Len()
}

func (lru *LruCache) Set(k, v interface{}) error {
	if lru.dataList == nil {
		return errors.New("lru is nil")
	}

	//map中存在，调整至链表头部
	if elem, ok := lru.dataMap[k]; ok {
		elem.Value.(*CacheNode).value = v
		lru.dataList.MoveToFront(elem)
		return nil
	}

	//调整至链表头部，并add to map
	newElem := lru.dataList.PushFront(&CacheNode{k, v})
	lru.dataMap[k] = newElem

	//超过capacity 移除最后一个
	if lru.dataList.Len() > lru.capacity {
		//移除最后一个
		lastElem := lru.dataList.Back()
		if lastElem == nil {
			//链表为空
			return nil
		}
		cacheNode := lastElem.Value.(*CacheNode)
		delete(lru.dataMap, cacheNode.key)
		lru.dataList.Remove(lastElem)
	}

	return nil
}

func (lru *LruCache) Get(k interface{}) (v interface{}, ret bool, err error) {
	if lru.dataList == nil {
		return v, false, errors.New("lru is nil")
	}

	//当前node调整至链表头部
	if elem, ok := lru.dataMap[k]; ok {
		lru.dataList.PushFront(elem)
		return elem.Value.(*CacheNode).value, true, nil
	}
	return v, false, nil
}

func (lru *LruCache) Remove(k interface{}) bool {
	if lru.dataList == nil {
		return false
	}

	if elem, ok := lru.dataMap[k]; ok {
		delete(lru.dataMap, k)
		lru.dataList.Remove(elem)
		return true
	}
	return false
}
