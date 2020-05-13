package caddyredis

import (
	"github.com/journeymidnight/yig-front-caddy"
)

type Config struct {
	Address  []string
	Password string
	RedisInfo
}

func NewConfig() *Config {
	cfg := new(Config)
	return cfg
}

func MakeRedisConfig(group *Config) *Redis {
	var redis *Redis
	redis = newRedis(*group)
	return redis
}

// ConfigGetter gets a Config keyed by key.
type ConfigGetter func(c *caddy.Controller) *Config

var configGetters = make(map[string]ConfigGetter)

func RegisterConfigGetter(serverType string, fn ConfigGetter) {
	configGetters[serverType] = fn
}
