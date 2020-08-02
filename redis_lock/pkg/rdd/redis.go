package rdd

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

//获取锁
func Lock(pool *redis.Pool, key string, ex uint) error {
	//set key value EX xx NX
	if setNxCmd, err := pool.Get().Do("SET", key, ex, "EX", ex, "NX"); err == nil {
		fmt.Println("send cmd = ", setNxCmd)
		//写入成功 锁持有时间ex(s)
		if setNxCmd != nil {
			//获得锁
			fmt.Println("get lock success")
			return nil
		} else {
			//锁冲突
			return fmt.Errorf("get lock failed: lock conflict")
		}
	} else {
		fmt.Printf("get lock err: %v", err)
	}
	return nil
}

//删除锁
func UnLockUnSafe(pool *redis.Pool, key string) bool {
	//直接删除锁key 会存在问题
	//如果删除之前 该key已经超时且被其他进程x获得锁 将会删除其他进程x的锁 锁被释放后
	//进程x删除不了自己的锁 又有其他进程再次获得锁 进而雪崩
	if delCmd, err := pool.Get().Do("DEL", key); err != nil {
		fmt.Println("del cmd = ", delCmd)
		//删除锁失败 会导致其他进程长时间等待锁
		return false
	}
	return true
}
