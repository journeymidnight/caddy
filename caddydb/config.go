package caddydb

import (
	"fmt"
	"github.com/journeymidnight/yig-front-caddy"
	"github.com/journeymidnight/yig-front-caddy/caddydb/clients/tidbclient"
	"strings"
)

type Config struct {
	Hostname string
	Clients  []string
	DBInfo   tidbclient.DBInfo
}

func NewConfig() *Config {
	cfg := new(Config)
	return cfg
}

func MakeDBConfig(cfg *Config) map[string]*tidbclient.TidbClient {
	clients := make(map[string]*tidbclient.TidbClient)
	for _, conn := range cfg.Clients {
		keyAll := strings.Split(conn, "/")
		key := keyAll[1]
		if clients[key] != nil {
			continue
		}
		clients[key] = tidbclient.NewTidbClient(conn, cfg.DBInfo)
	}
	var keys []string
	for key, _ := range clients {
		keys = append(keys, key)
	}
	fmt.Println("Already loaded database connections:", keys)
	return clients
}

// ConfigGetter gets a Config keyed by key.
type ConfigGetter func(c *caddy.Controller) *Config

var configGetters = make(map[string]ConfigGetter)

func RegisterConfigGetter(serverType string, fn ConfigGetter) {
	configGetters[serverType] = fn
}
