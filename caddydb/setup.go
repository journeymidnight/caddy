package caddydb

import (
	"fmt"
	caddy "github.com/journeymidnight/yig-front-caddy"
	"strconv"
)

func init() {
	caddy.RegisterPlugin("database", caddy.Plugin{Action: setupDB})
}

func setupDB(c *caddy.Controller) error {
	configGetter, ok := configGetters[c.ServerType()]
	if !ok {
		return fmt.Errorf("no caddydatabase.ConfigGetter for %s server type; must call RegisterConfigGetter", c.ServerType())
	}
	config := configGetter(c)
	if config == nil {
		return fmt.Errorf("no caddydatabase.Config to set up for %s", c.Key)
	}

	var hadBlock bool
	for c.Next() {
		for c.NextBlock() {
			hadBlock = true
			switch c.Val() {
			case "databases":
				args := c.RemainingArgs()
				if len(args) <= 0 {
					return c.Errf("Yig database does not match: '%s'", c.ArgErr())
				}
				config.Clients = args
			case "db_max_idle_conns":
				if !c.NextArg() {
					return c.Errf("Wrong key database max idle connections: '%s'", c.ArgErr())
				}
				int, err := strconv.Atoi(c.Val())
				if err != nil {
					return c.Errf("Wrong key database max idle connections: '%s'", err)
				}
				config.DBInfo.DBMaxIdleConns = int
			case "db_max_open_conns":
				if !c.NextArg() {
					return c.Errf("Wrong key database max open connections: '%s'", c.ArgErr())
				}
				int, err := strconv.Atoi(c.Val())
				if err != nil {
					return c.Errf("Wrong key database max open connections: '%s'", err)
				}
				config.DBInfo.DBMaxOpenConns = int
			case "db_conn_max_life_seconds":
				if !c.NextArg() {
					return c.Errf("Wrong key database max connection life second: '%s'", c.ArgErr())
				}
				int, err := strconv.Atoi(c.Val())
				if err != nil {
					return c.Errf("Wrong key database max connection life second: '%s'", err)
				}
				config.DBInfo.DBConnMaxLifeSeconds = int
			default:
				return c.Errf("Unknown subdirective '%s'", c.Val())
			}
		}
		//If the configuration block is enabled, it must contain parameters
		if !hadBlock {
			return c.ArgErr()
		}
	}
	return nil
}
