package v1

import (
	"testing"
)

func TestLruV1(t *testing.T) {
	lru := NewLruCacheV1(3)

	//set node
	lru.Set(10, "value10")
	lru.Set(20, "value20")
	lru.Set(30, "value30")
	lru.Set(10, "value10")
	lru.Set(50, "value50")

	t.Log("size :", lru.Size())
	for v := range lru.dataMap {
		t.Log("v = ", v)
	}

	if lru.Remove(30) {
		t.Log("Remove(30) : true ")
	} else {
		t.Log("Remove(30) : false ")
	}
	t.Log("LRU Size:", lru.Size())

}
