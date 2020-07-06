package test

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"go-insight/redis_lock/pkg"
	"testing"
)

func NewRedisClient() redis.Conn {
	redisClient, err := redis.Dial("tcp", pkg.RedisAddr)
	if err != nil {
		fmt.Printf("conn to redis err: %v", err)
	}
	defer redisClient.Close()
	return redisClient
}

func TestRddRedisLock(t *testing.T) {
	//rds := NewRedisClient()

}
