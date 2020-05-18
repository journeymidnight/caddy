package caddydb

import (
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
	var yigC, caddyC string
	yigC = DefaultYigSource
	caddyC = DefaultCaddySource
	clients := []string{yigC, caddyC}
	var dbInfo tidbclient.DBInfo
	dbInfo.DBMaxIdleConns = DefaultDBMaxIdleConns
	dbInfo.DBMaxOpenConns = DefaultDBMaxOpenConns
	dbInfo.DBConnMaxLifeSeconds = DefaultDBConnMaxLifeSeconds
	cfg := new(Config)
	cfg.DBInfo = dbInfo
	cfg.Clients = clients
	return cfg
}

func MakeDBConfig(group []*Config) map[string]*tidbclient.TidbClient {
	clients := make(map[string]*tidbclient.TidbClient)
	for _, cfg := range group {
		for _, conn := range cfg.Clients {
			keyAll := strings.Split(conn, "/")
			key := keyAll[1]
			if clients[key] != nil {
				continue
			}
			clients[key] = tidbclient.NewTidbClient(conn, cfg.DBInfo)
		}
	}
	return clients
}

// ConfigGetter gets a Config keyed by key.
type ConfigGetter func(c *caddy.Controller) *Config

var configGetters = make(map[string]ConfigGetter)

func RegisterConfigGetter(serverType string, fn ConfigGetter) {
	configGetters[serverType] = fn
}

const (
	DefaultYigSource            = "root:@tcp(10.5.0.17:4000)/yig"
	DefaultCaddySource          = "root:@tcp(10.5.0.17:4000)/caddy"
	DefaultDBMaxIdleConns       = 1024
	DefaultDBMaxOpenConns       = 10240
	DefaultDBConnMaxLifeSeconds = 300
	ClientsCacheInstStorageKey  = "cli_cert_cache"
)
