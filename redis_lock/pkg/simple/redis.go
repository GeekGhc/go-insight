package simple

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"go-insight/redis_lock/pkg"
	"time"
)

func GetLock(lockKey string, ex uint, retry int) error {
	if retry <= 0 {
		retry = 10
	}
	redisClient, err := redis.Dial("tcp", pkg.RedisAddr)
	if err != nil {
		fmt.Printf("conn to redis err: %v", err)
	}
	defer redisClient.Close()
	//as random value
	ts := time.Now()
	for i := 1; i <= retry; i++ {
		//sleep if not first time
		if i > 1 {
			time.Sleep(time.Second)
		}

		v, err := redisClient.Do("SET", lockKey, ts, "EX", ex, "NX")
		if err == nil {
			if v == nil {
				//锁存在 互斥
				fmt.Println("get lock failed,retry time: ", i)
			} else {
				//获得锁
				fmt.Println("get lock success")
				break
			}
		} else {
			fmt.Printf("get lock err: %v", err)
		}

		//到达最大重试次数
		if i >= retry {
			err = fmt.Errorf("get lock failed with max retry times\n")
			return err
		}
	}
	return nil
}

func UnLock(lockKey string) error {
	redisClient, err := redis.Dial("tcp", pkg.RedisAddr)
	if err != nil {
		fmt.Printf("conn to redis err: %v", err)
	}
	defer redisClient.Close()

	v, err := redis.Bool(redisClient.Do("DEL", lockKey))
	if err == nil {
		if v {
			fmt.Println("unlock success")
		} else {
			fmt.Println("unlock failed")
		}
	} else {
		fmt.Printf("unLock err: %v", err)
	}
	return nil
}
