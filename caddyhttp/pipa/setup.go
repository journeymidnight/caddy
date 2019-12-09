package pipa

import (
	caddy "github.com/journeymidnight/yig-front-caddy"
	"github.com/journeymidnight/yig-front-caddy/caddyhttp/httpserver"
	"strconv"
)

func init() {
	caddy.RegisterPlugin("pipa", caddy.Plugin{
		ServerType: "http",
		Action:     setup,
	})
}

// setup configures a new mime middleware instance.
func setup(c *caddy.Controller) error {
	redisMaxIdle, redisAddress, redisPwd, err := pipaParse(c)
	if err != nil {
		return err
	}
	httpserver.GetConfig(c).AddMiddleware(func(next httpserver.Handler) httpserver.Handler {
		return Pipa{Next: next, redis: newRedisPool(redisMaxIdle, redisAddress, redisPwd)}
	})
	return nil
}

func pipaParse(c *caddy.Controller) (redisMaxIdle int, redisAddress, redisPwd string, err error) {
	for c.Next() {
		for c.NextBlock() {
			ext := c.Val()
			switch ext {
			case "redis_maxIdle":
				c.NextArg()
				key := c.Val()
				redisMaxIdle, err = strconv.Atoi(key)
				if err != nil {
					return -1, "", "", err
				}
				break
			case "redis_address":
				c.NextArg()
				redisAddress = c.Val()
				break
			case "redis_password":
				c.NextArg()
				redisPwd = c.Val()
				break
			}
		}
	}
	return
}
