package test

import (
	"fmt"
	"go-insight/redis_lock/pkg/simple"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestSimpleRedisLock(t *testing.T) {
	var wg sync.WaitGroup

	key := "simple_redis_key"

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(pId int) {
			defer wg.Done()
			time.Sleep(time.Second)
			//getLock  每个协程最多持有锁5秒  最多重试10次
			err := simple.GetLock(key, 5, 10)
			if err != nil {
				t.Errorf("routine %d get lock failed: %v", pId, err)
			}
			//sleep for random 模拟持锁的操作事件时间(不操作超时时间5s)
			rand.Seed(time.Now().Unix())
			randomNum := time.Duration(rand.Intn(5)) + 1
			fmt.Println("random time: ", randomNum*time.Second)
			time.Sleep(randomNum * time.Second)
			//unlock
			err = simple.UnLock(key)
			if err != nil {
				t.Errorf("routine %d unlock failed: %v", pId, err)
			}
		}(i)
	}
	wg.Wait()
	t.Log("done...")
}
