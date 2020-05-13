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
	secretKey, reservedOrigins, err := pipaParse(c)
	if err != nil {
		return err
	}
	httpserver.GetConfig(c).AddMiddleware(func(next httpserver.Handler) httpserver.Handler {
		return Pipa{
			SecretKey:       secretKey,
			ReservedOrigins: reservedOrigins,
		}
	})
	return nil
}

func pipaParse(c *caddy.Controller) (secretKey string, reservedOrigins []string, err error) {
	for c.Next() {
		for c.NextBlock() {
			ext := c.Val()
			switch ext {
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
