package rdd

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"time"
)

func expiredTime(second uint) (int64, int64) {
	now := time.Now()
	//秒=>纳秒
	return now.UnixNano(), now.Add(time.Duration(second * 1000000000)).UnixNano()
}

//获取锁
func Lock(pool *redis.Pool, key string, ts uint) (error, int64) {
	rc := pool.Get()
	_, ex := expiredTime(ts)
	//set key value EX xx NX
	if setNxCmd, err := rc.Do("SET", key, ex, "EX", ts, "NX"); err == nil {
		//写入成功 锁持有时间ex(s)
		if setNxCmd != nil {
			//获得锁
			log.Println("LockRetry get lock success")
			return nil, ex
		} else {
			//锁冲突
			return fmt.Errorf("LockRetry get lock failed: lock conflict"), 0
		}
	} else {
		log.Printf("LockRetry get lock err: %v", err)
	}
	return nil, 0
}

//删除锁(UnSafe)
func UnLockUnSafe(pool *redis.Pool, key string) bool {
	//直接删除锁key 会存在问题
	//如果删除之前 该key已经超时且被其他进程x获得锁 将会删除其他进程x的锁 锁被释放后
	//进程x删除不了自己的锁 又有其他进程再次获得锁 进而雪崩
	rc := pool.Get()
	if _, err := rc.Do("DEL", key); err != nil {
		log.Printf("UnLockUnSafe del cmd err: %v", err)
		//删除锁失败 会导致其他进程长时间等待锁
		return false
	}
	return true
}

//删除锁(Safe)
func UnLockSafe(pool *redis.Pool, key string, ts int64) bool {
	//直接删除key 不过增加安全事件
	//若在安全时间内 key即将过期，就等待key过期(其他进程则判断超时) 防止雪崩
	rc := pool.Get()
	getCmd, err := rc.Do("GET", key)
	val, err := redis.Uint64(getCmd, err)
	if err != nil {
		log.Printf("UnLockSafe get cmd err: %v", err)
		return false
	}

	//安全时间判断
	now := time.Now().UnixNano()
	if now+ts*1000000000 > int64(val) {
		log.Println("UnLockSafe the key is going to expire.")
		return false
	}

	if _, err := rc.Do("DEL", key); err != nil {
		//删除锁失败 会导致其他进程长时间等待锁
		log.Printf("UnLockSafe del cmd err: %v", err)
		return false
	}
	return true
}

func UnLock(pool *redis.Pool, key string, ex int64) bool {
	rc := pool.Get()
	if ex <= 0 {
		//没有获得锁 不能解锁
		return false
	}

	getCmd, err := rc.Do("GET", key)
	val, err := redis.Uint64(getCmd, err)
	if err != nil {
		log.Printf("UnLock get cmd err: %v", err)
		return false
	}

	//进程持有的锁匹配-进程只能删除自己的锁
	if int64(val) == ex {
		//开始解锁
		if _, err := rc.Do("DEL", key); err != nil {
			log.Printf("UnLock del cmd err: %v", err)
			//删除锁失败
			return false
		}
		return true
	}
	return false
}

func LockRetry(pool *redis.Pool, key string, ts uint, retryTimes int) (bool, int64) {
	rc := pool.Get()
	for i := 0; i < retryTimes; i++ {
		_, ex := expiredTime(ts)
		//set key value EX xx NX
		if setNxCmd, err := rc.Do("SET", key, ex, "EX", ts, "NX"); err == nil {
			//写入成功 锁持有时间ex(s)
			if setNxCmd != nil {
				//获得锁
				log.Println("LockRetry get lock success")
				return true, ex
			} else {
				//锁冲突
				log.Print("LockRetry get lock failed: lock conflict.")
				continue
			}
		} else {
			log.Printf("LockRetry get lock err: %v", err)
			continue
		}
	}
	return false, 0
}

func SyncDo(pool *redis.Pool, key string, ts uint, do func(timeout chan bool) chan bool) error {
	ticker := time.NewTicker(time.Duration(ts * 1000000000))
	locked, ex := LockRetry(pool, key, ts, 10)
	if !locked {
		return fmt.Errorf("SyncDo lock the key %s failed", key)
	}

	defer func() {
		unlocked := UnLock(pool, key, ex)
		if !unlocked {
			UnLockUnSafe(pool, key)
		}
	}()

	timeout := make(chan bool, 1)
	doRet := do(timeout)
	select {
	case <-ticker.C:
		timeout <- true
		return fmt.Errorf("timeout")
	case <-doRet:
		return nil
	}
}
