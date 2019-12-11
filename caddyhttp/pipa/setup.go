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
	redisMaxIdle, redisAddress, redisPwd, secretKey, err := pipaParse(c)
	if err != nil {
		return err
	}
	httpserver.GetConfig(c).AddMiddleware(func(next httpserver.Handler) httpserver.Handler {
		return Pipa{Next: next, redis: newRedisPool(redisMaxIdle, redisAddress, redisPwd), SecretKey: secretKey}
	})
	return nil
}

func pipaParse(c *caddy.Controller) (redisMaxIdle int, redisAddress, redisPwd, secretKey string, err error) {
	for c.Next() {
		for c.NextBlock() {
			ext := c.Val()
			switch ext {
			case "redis_maxIdle":
				if !c.NextArg() {
					err = c.ArgErr()
					return
				}
				key := c.Val()
				redisMaxIdle, err = strconv.Atoi(key)
				if err != nil {
					return
				}
				break
			case "redis_address":
				if !c.NextArg() {
					err = c.ArgErr()
					return
				}
				redisAddress = c.Val()
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
			}
		}
	}
	return
}
