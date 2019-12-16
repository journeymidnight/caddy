package customdomain

import (
	"github.com/journeymidnight/yig-front-caddy"
	"github.com/journeymidnight/yig-front-caddy/caddyhttp/httpserver"
	"time"
	"zvelo.io/ttlru"
)

func init() {
	caddy.RegisterPlugin("customdomain", caddy.Plugin{
		ServerType: "http",
		Action:     setup,
	})
}

// setup configures a new host middleware instance.
func setup(c *caddy.Controller) error {
	customDomainFlag, secretKey, sealKey, err := hostParse(c)
	if err != nil {
		return err
	}
	httpserver.GetConfig(c).AddMiddleware(func(next httpserver.Handler) httpserver.Handler {
		return Domain{
			Next:             next,
			CustomDomainFlag: customDomainFlag,
			SecretKey:        secretKey,
			SealKey:          sealKey,
			Cache:            ttlru.New(1024, ttlru.WithTTL(10*time.Minute)),
		}
	})
	return nil
}

func hostParse(c *caddy.Controller) (customDomainFlag string, secretKey string, sealKey string, err error) {
	for c.Next() {
		for c.NextBlock() {
			ext := c.Val()
			switch ext {
			case "custom_domainflag":
				if !c.NextArg() {
					return "", "", "", c.ArgErr()
				}
				customDomainFlag = c.Val()
				break
			case "secret_key":
				if !c.NextArg() {
					return "", "", "", c.ArgErr()
				}
				secretKey = c.Val()
			case "seal_key":
				if !c.NextArg() {
					return "", "", "", c.ArgErr()
				}
				sealKey = c.Val()
			default:
				return "", "", "", c.ArgErr()
			}
		}
	}
	return
}
