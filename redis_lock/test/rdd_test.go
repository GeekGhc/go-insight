package test

import (
	"flag"
	"github.com/gomodule/redigo/redis"
	"go-insight/redis_lock/pkg/rdd"
	"testing"
	"time"
)

var (
	pool        *redis.Pool
	redisServer = flag.String("127.0.0.1", ":6379", "")
)

func newPool(addr string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		// Dial or DialContext must be set. When both are set, DialContext takes precedence over Dial.
		Dial: func() (redis.Conn, error) { return redis.Dial("tcp", addr) },
	}
}

func TestRddRedisLock(t *testing.T) {
	pool = newPool(*redisServer)
	key := "rdd_redis_key"

	//获取锁
	if err := rdd.Lock(pool, key, 20); err != nil {
		t.Error("get lock failed...")
	}

}
