package model

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"time"
)

var pool *redis.Pool

func init()  {
	pool = newRedisPool()
}

func newRedisPool() *redis.Pool {
	redisHost := "127.0.0.1:6379"
	redisPass := "QWER1234"

	return &redis.Pool{
		MaxIdle:     50,
		MaxActive:   30,
		IdleTimeout: 300 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", redisHost, redis.DialPassword(redisPass))
			if err != nil {
				fmt.Println(err)
				return nil, err
			}
			//2、访问认证
			if _, err = c.Do("AUTH", redisPass); err != nil {
				c.Close()
				return nil, err
			}
			c.Do("SELECT", "0")
			return c, nil
		},
		//定时检查redis是否出状况
		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := conn.Do("PING")
			return err
		},
	}
}

func RedisPool() *redis.Pool {
	return pool
}

