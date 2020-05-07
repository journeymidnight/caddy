package pipa

import (
	"github.com/go-redis/redis/v7"
	. "github.com/journeymidnight/yig-front-caddy/caddyerrors"
	"strings"
	"time"
)

type Redis interface {
	pushRequest(data []byte) (err error)
	popResponse(data TaskData) (result []byte, err error)
	getImageFromRedis(url string) (result []byte, err error)
}

func newRedis(server []string, password string) (redisConn Redis) {
	if len(server) == 1 {
		r := InitializeSingle(server[0], password)
		redisConn = r.(Redis)
		return redisConn
	} else {
		r := InitializeCluster(server, password)
		redisConn = r.(Redis)
		return redisConn
	}
}

type SingleRedis struct {
	client *redis.Client
}

func InitializeSingle(server, password string) interface{} {
	options := &redis.Options{
		Addr:         server,
		DialTimeout:  time.Duration(5) * time.Second,
		ReadTimeout:  time.Duration(5) * time.Second,
		WriteTimeout: time.Duration(5) * time.Second,
		IdleTimeout:  time.Duration(30) * time.Second,
	}
	if password != "" {
		options.Password = password
	}
	client := redis.NewClient(options)
	_, err := client.Ping().Result()
	if err != nil {
		PIPA.Log.Error("redis PING err:", err)
		return nil
	}
	r := &SingleRedis{client: client}
	return interface{}(r)
}

func (s *SingleRedis) pushRequest(data []byte) (err error) {
	redis_conn := s.client.Conn()
	defer redis_conn.Close()
	_, err = redis_conn.LPush("taskQueue", data).Result()
	return
}

func (s *SingleRedis) popResponse(data TaskData) (result []byte, err error) {
	redis_conn := s.client.Conn()
	defer redis_conn.Close()
	response, err := redis_conn.BRPop(time.Duration(30)*time.Second, data.Uuid).Result()
	if response == nil {
		return nil, ErrTimeout
	}
	pipaResult := strings.Split(response[1], ",")
	if pipaResult[0] != "200" {
		return []byte(response[1]), ErrInternalServer
	} else {
		result, err = redis_conn.Get(data.Url).Bytes()
		if err != nil {
			return
		}
	}
	return
}

func (s *SingleRedis) getImageFromRedis(url string) (result []byte, err error) {
	redis_conn := s.client.Conn()
	defer redis_conn.Close()
	result, err = redis_conn.Get(url).Bytes()
	if err != nil {
		return
	}
	PIPA.Log.Info("Success get image from redis directly!")
	return
}

type ClusterRedis struct {
	cluster *redis.ClusterClient
}

func InitializeCluster(server []string, password string) interface{} {
	clusterRedis := &redis.ClusterOptions{
		Addrs:        server,
		DialTimeout:  time.Duration(5) * time.Second,
		ReadTimeout:  time.Duration(5) * time.Second,
		WriteTimeout: time.Duration(5) * time.Second,
		IdleTimeout:  time.Duration(30) * time.Second,
	}
	if password != "" {
		clusterRedis.Password = password
	}
	cluster := redis.NewClusterClient(clusterRedis)
	_, err := cluster.Ping().Result()
	if err != nil {
		PIPA.Log.Error("Cluster Mode redis PING err:", err)
		return nil
	}
	r := &ClusterRedis{cluster: cluster}
	return interface{}(r)
}

func (c *ClusterRedis) pushRequest(data []byte) (err error) {
	redis_conn := c.cluster
	_, err = redis_conn.LPush("taskQueue", data).Result()
	return
}

func (c *ClusterRedis) popResponse(data TaskData) (result []byte, err error) {
	redis_conn := c.cluster
	response, err := redis_conn.BRPop(time.Duration(30)*time.Second, data.Uuid).Result()
	if response == nil {
		return nil, ErrTimeout
	}
	pipaResult := strings.Split(response[1], ",")
	if pipaResult[0] != "200" {
		return []byte(response[1]), ErrInternalServer
	} else {
		result, err = redis_conn.Get(data.Url).Bytes()
		if err != nil {
			return
		}
	}
	return
}

func (c *ClusterRedis) getImageFromRedis(url string) (result []byte, err error) {
	redis_conn := c.cluster
	result, err = redis_conn.Get(url).Bytes()
	if err != nil {
		return
	}
	PIPA.Log.Info("Success get image from redis directly!")
	return
}
