package pipa

import (
	caddy "github.com/journeymidnight/yig-front-caddy"
	"github.com/journeymidnight/yig-front-caddy/caddyhttp/httpserver"
)

func init() {
	caddy.RegisterPlugin("pipa", caddy.Plugin{
		ServerType: "http",
		Action:     setup,
	})
}

// setup configures a new mime middleware instance.
func setup(c *caddy.Controller) error {
	redisAddrs, redisPwd, secretKey, reservedOrigins, err := pipaParse(c)
	if err != nil {
		return err
	}
	httpserver.GetConfig(c).AddMiddleware(func(next httpserver.Handler) httpserver.Handler {
		return Pipa{
			Next:            next,
			redis:           newRedis(redisAddrs, redisPwd),
			SecretKey:       secretKey,
			ReservedOrigins: reservedOrigins,
		}
	})
	return nil
}

func pipaParse(c *caddy.Controller) (redisAddress []string, redisPwd, secretKey string, reservedOrigins []string, err error) {
	for c.Next() {
		for c.NextBlock() {
			ext := c.Val()
			switch ext {
			case "redis_address":
				if !c.NextArg() {
					err = c.ArgErr()
					return
				}
				redisAddress = append(redisAddress, c.Val())
				redisAddress = getAddrs(c, redisAddress)
				break
			case "redis_password":
				if !c.NextArg() {
					err = c.ArgErr()
					return
				}
				redisPwd = c.Val()
				break
			case "secret_key":
				if !c.NextArg() {
					err = c.ArgErr()
					return
				}
				secretKey = c.Val()
			case "reserved_origins":
				args := c.RemainingArgs()
				if len(args) <= 0 {
					err = c.ArgErr()
					return
				}

			}
		}
	}
	return
}

func getAddrs(c *caddy.Controller, addrs []string) []string {
	if c.NextArg() {
		addrs = append(addrs, c.Val())
		addrs = getAddrs(c, addrs)
	}
	return addrs
}
