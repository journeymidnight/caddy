package lib

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
)

func RedisConn() (redis.Conn, error) {
	pwd := redis.DialPassword(RedisPwd)
	red, err := redis.Dial("tcp", RedisDir, pwd)
	if err != nil {
		fmt.Println("Redis connection", red, err)
		return nil, err
	}
	_, err = red.Do("PING")
	if err != nil {
		return nil, err
	}
	return red, nil
}

func RedisFlushAll()  {
	red, err := RedisConn()
	if err != nil {
		panic(err)
	}
	_, err = red.Do("FLUSHALL")
	if err != nil {
		panic(err)
	}
	return
}

type FinishTask struct {
	code int
	uuid string
	url  string
	blob []byte
}

type TaskData struct {
	Uuid string `json:"uuid"`
	Url  string `json:"url"`
}

func RedisPipa() {
	red, err := RedisConn()
	if err != nil {
		panic(err)
	}
	response, err := redis.Strings(red.Do("BRPOP", "taskQueue", 30))
	if response == nil {
		panic("response is nil")
	}
	data := &TaskData{}
	err = json.Unmarshal([]byte(response[1]), data)
	if err != nil {
		panic(err)
	}
	res := &FinishTask{}
	res.uuid = data.Uuid
	res.url = data.Url
	res.code = 200
	res.blob = []byte("hehehehehehe")
	_, err = red.Do("MULTI")
	if err != nil {
		panic(err)
	}
	_, err = red.Do("SET", res.url, res.blob)
	if err != nil {
		panic(err)
	}
	_, err = red.Do("LPUSH", res.uuid, res.code)
	if err != nil {
		panic(err)
	}
	_, err = red.Do("EXEC")
	if err != nil {
		panic(err)
	}
}
