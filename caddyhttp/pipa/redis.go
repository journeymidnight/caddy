package pipa

import (
	"github.com/garyburd/redigo/redis"
	"github.com/journeymidnight/yig-front-caddy/caddyerrors"
	"time"
)

func newRedisPool(MaxIdle int, server, password string) *redis.Pool {
	pwd := redis.DialPassword(password)
	return &redis.Pool{
		MaxIdle:     MaxIdle,
		IdleTimeout: 60 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server, pwd)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func pushRequest(data []byte, pool *redis.Pool) (err error) {
	redis_conn := pool.Get()
	_, err = redis_conn.Do("LPUSH", "taskQueue", data)
	err = redis_conn.Close()
	return
}

func popResponse(data TaskData, pool *redis.Pool) (result []byte, err error) {
	redis_con := pool.Get()
	defer redis_con.Close()
	response, err := redis.Strings(redis_con.Do("BRPOP", data.Uuid, 30))
	if response == nil {
		return nil, caddyerrors.ErrTimeout
	}
	if response[1] != "200" {
		return nil, caddyerrors.ErrInternalServer
	} else {
		result, err = redis.Bytes(redis_con.Do("GET", data.Url))
		if err != nil {
			return
		}
	}
	return
}
