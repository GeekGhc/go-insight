package rdd

import (
	"github.com/gomodule/redigo/redis"
	"strconv"
)

//获取锁
func Lock(rds redis.Conn, key string, timeoutMs uint) (bool, string) {
	if setNxCmd, err := rds.Do("SET", key, "EX", timeoutMs, "NX"); err == nil {
		//写入成功 锁持有时间ex
		if setNxCmd != nil {
			//获得锁
			return true, strconv.FormatInt(int64(timeoutMs), 10)
		}
	}
	return false, "0"
}
