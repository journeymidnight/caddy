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
	cfg.Address = []string{DEFAULTREDISADDR}
	cfg.Password = DEFAULTREDISPWD
	cfg.MaxRetries = DEFAULTREDISMAXRETRIES
	cfg.ConnectTimeout = DEFAULTREDISCONNTIMEOUT
	cfg.ReadTimeout = DEFAULTREDISREADTIMEOUT
	cfg.WriteTimeout = DEFAULTREDISWRITETIMEOUT
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

const (
	DEFAULTREDISADDR         = "redis:6379"
	DEFAULTREDISPWD          = "hehehehe"
	DEFAULTREDISMAXRETRIES   = 20
	DEFAULTREDISCONNTIMEOUT  = 5
	DEFAULTREDISREADTIMEOUT  = 5
	DEFAULTREDISWRITETIMEOUT = 5
)
