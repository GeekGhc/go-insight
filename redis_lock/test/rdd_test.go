package test

import (
	"flag"
	"github.com/gomodule/redigo/redis"
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
}
