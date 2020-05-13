package caddyredis

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	. "github.com/journeymidnight/yig-front-caddy/caddyerrors"
	"strings"
	"time"
)

type RedisInfo struct {
	MaxRetries     int
	ConnectTimeout int
	ReadTimeout    int
	WriteTimeout   int
}

type Redis struct {
	single  bool
	client  *SingleRedis
	cluster *ClusterRedis
}

var redisConn Redis

func (r *Redis) PushRequest(data []byte) (err error) {
	if r.single {
		return r.client.pushRequest(data)
	}
	return r.cluster.pushRequest(data)
}

func (r *Redis) PopResponse(uuid, url string) (result []byte, err error) {
	if r.single {
		return r.client.popResponse(uuid, url)
	}
	return r.cluster.popResponse(uuid, url)
}

func (r *Redis) GetImageFromRedis(url string) (result []byte, err error) {
	if r.single {
		return r.client.getImageFromRedis(url)
	}
	return r.cluster.getImageFromRedis(url)
}

func newRedis(info Config) *Redis {
	redis := &Redis{}
	fmt.Print("Redis is configured as:", info.Address, len(info.Address) == 1, "\n")
	if len(info.Address) == 1 {
		redis.single = true
		redis.client = InitializeSingle(info)
		return redis
	} else {
		redis.single = false
		redis.cluster = InitializeCluster(info)
		return redis
	}
}

type SingleRedis struct {
	client *redis.Client
}

func InitializeSingle(info Config) *SingleRedis {
	options := &redis.Options{
		Addr:         info.Address[0],
		DialTimeout:  time.Duration(5) * time.Second,
		ReadTimeout:  time.Duration(5) * time.Second,
		WriteTimeout: time.Duration(5) * time.Second,
		IdleTimeout:  time.Duration(30) * time.Second,
	}
	if info.Password != "" {
		options.Password = info.Password
	}
	client := redis.NewClient(options)
	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}
	return &SingleRedis{client: client}
}

func (s *SingleRedis) pushRequest(data []byte) (err error) {
	redis_conn := s.client.Conn()
	defer redis_conn.Close()
	_, err = redis_conn.LPush("taskQueue", data).Result()
	return
}

func (s *SingleRedis) popResponse(uuid, url string) (result []byte, err error) {
	redis_conn := s.client.Conn()
	defer redis_conn.Close()
	response, err := redis_conn.BRPop(time.Duration(30)*time.Second, uuid).Result()
	if response == nil {
		return nil, ErrTimeout
	}
	pipaResult := strings.Split(response[1], ",")
	if pipaResult[0] != "200" {
		return []byte(response[1]), ErrInternalServer
	} else {
		result, err = redis_conn.Get(url).Bytes()
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
	return
}

type ClusterRedis struct {
	cluster *redis.ClusterClient
}

func InitializeCluster(info Config) *ClusterRedis {
	clusterRedis := &redis.ClusterOptions{
		Addrs:        info.Address,
		DialTimeout:  time.Duration(info.ConnectTimeout) * time.Second,
		ReadTimeout:  time.Duration(info.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(info.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(30) * time.Second,
	}
	if info.Password != "" {
		clusterRedis.Password = info.Password
	}
	cluster := redis.NewClusterClient(clusterRedis)
	_, err := cluster.Ping().Result()
	if err != nil {
		panic(err)
	}
	return &ClusterRedis{cluster: cluster}
}

func (c *ClusterRedis) pushRequest(data []byte) (err error) {
	redis_conn := c.cluster
	_, err = redis_conn.Ping().Result()
	if err != nil {
		for {
			if _, err = redis_conn.Ping().Result(); err == nil {
				break
			}
		}
	}
	_, err = redis_conn.LPush("taskQueue", data).Result()
	if err != nil {
		fmt.Println("err LPush info:", err)
	}
	return
}

func (c *ClusterRedis) popResponse(uuid, url string) (result []byte, err error) {
	redis_conn := c.cluster
	response, err := redis_conn.BRPop(time.Duration(30)*time.Second, uuid).Result()
	if response == nil {
		return nil, ErrTimeout
	}
	pipaResult := strings.Split(response[1], ",")
	if pipaResult[0] != "200" {
		return []byte(response[1]), ErrInternalServer
	} else {
		result, err = redis_conn.Get(url).Bytes()
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
	return
}
