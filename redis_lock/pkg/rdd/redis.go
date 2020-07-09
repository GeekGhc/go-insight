package rdd

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

//获取锁
func Lock(pool *redis.Pool, key string, timeoutMs uint) error {
	if setNxCmd, err := pool.Get().Do("SET", key, "EX", timeoutMs, "NX"); err == nil {
		//写入成功 锁持有时间ex
		if setNxCmd != nil {
			//获得锁
			fmt.Println("get lock success")
			return nil
		} else {
			//锁冲突
			return fmt.Errorf("get lock failed: lock conflict")
		}
	}
	return nil
}
