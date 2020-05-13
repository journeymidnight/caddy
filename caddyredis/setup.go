package caddyredis

import (
	"fmt"
	caddy "github.com/journeymidnight/yig-front-caddy"
	"strconv"
)

func init() {
	caddy.RegisterPlugin("redis", caddy.Plugin{Action: setupRedis})
}

func setupRedis(c *caddy.Controller) error {
	configGetter, ok := configGetters[c.ServerType()]
	if !ok {
		return fmt.Errorf("no caddytls.ConfigGetter for %s server type; must call RegisterConfigGetter", c.ServerType())
	}
	config := configGetter(c)
	if config == nil {
		return fmt.Errorf("no caddytls.Config to set up for %s", c.Key)
	}

	var hadBlock bool
	for c.Next() {
		for c.NextBlock() {
			hadBlock = true
			switch c.Val() {
			case "address":
				args := c.RemainingArgs()
				if len(args) <= 0 {
					return c.Errf("Redis does not match: '%s'", c.ArgErr())
				}
				config.Address = args
			case "password":
				if !c.NextArg() {
					return c.Errf("Wrong key redis password: '%s'", c.ArgErr())
				}
			case "max_retries":
				if !c.NextArg() {
					return c.Errf("Wrong key redis max retries: '%s'", c.ArgErr())
				}
				int, err := strconv.Atoi(c.Val())
				if err != nil {
					return c.Errf("Wrong key redis max retries: '%s'", err)
				}
				config.MaxRetries = int
			case "conn_timeout":
				if !c.NextArg() {
					return c.Errf("Wrong key redis connection timeout: '%s'", c.ArgErr())
				}
				int, err := strconv.Atoi(c.Val())
				if err != nil {
					return c.Errf("Wrong key redis connection timeout: '%s'", err)
				}
				config.ConnectTimeout = int
			case "read_timeout":
				if !c.NextArg() {
					return c.Errf("Wrong key redis read timeout: '%s'", c.ArgErr())
				}
				int, err := strconv.Atoi(c.Val())
				if err != nil {
					return c.Errf("Wrong key redis read timeout: '%s'", err)
				}
				config.ReadTimeout = int
			case "write_timeout":
				if !c.NextArg() {
					return c.Errf("Wrong key redis write timeout: '%s'", c.ArgErr())
				}
				int, err := strconv.Atoi(c.Val())
				if err != nil {
					return c.Errf("Wrong key redis write timeout: '%s'", err)
				}
				config.ReadTimeout = int
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
