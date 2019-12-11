package pipa

import (
	caddy "github.com/journeymidnight/yig-front-caddy"
	"github.com/journeymidnight/yig-front-caddy/caddyhttp/clients/clients/tidbclient"
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
	redisMaxIdle, redisAddress, redisPwd, secretKey, s3Source, caddySource, db, err := pipaParse(c)
	if err != nil {
		return err
	}
	httpserver.GetConfig(c).AddMiddleware(func(next httpserver.Handler) httpserver.Handler {
		return Pipa{
			Next:      next,
			redis:     newRedisPool(redisMaxIdle, redisAddress, redisPwd),
			Client:    tidbclient.NewCustomDomainClient(s3Source, caddySource, db),
			SecretKey: secretKey,
		}
	})
	return nil
}

func pipaParse(c *caddy.Controller) (redisMaxIdle int, redisAddress, redisPwd, secretKey, s3Source, caddySource string, db tidbclient.DBInfo, err error) {
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
			case "s3_db":
				if !c.NextArg() {
					err = c.ArgErr()
					return
				}
				s3Source = c.Val()
				break
			case "caddy_db":
				if !c.NextArg() {
					err = c.ArgErr()
					return
				}
				caddySource = c.Val()
				break
			case "db_max_idle_conns":
				if !c.NextArg() {
					err = c.ArgErr()
					return
				}
				int, err := strconv.Atoi(c.Val())
				if err != nil {
					return -1, "", "", "", "", "", db, err
				}
				db.DBMaxIdleConns = int
			case "db_max_open_conns":
				if !c.NextArg() {
					err = c.ArgErr()
					return
				}
				int, err := strconv.Atoi(c.Val())
				if err != nil {
					return -1, "", "", "", "", "", db, err
				}
				db.DBMaxOpenConns = int
			case "db_conn_max_life_seconds":
				if !c.NextArg() {
					err = c.ArgErr()
					return
				}
				int, err := strconv.Atoi(c.Val())
				if err != nil {
					return -1, "", "", "", "", "", db, err
				}
				db.DBConnMaxLifeSeconds = int
			}
		}
	}
	return
}
