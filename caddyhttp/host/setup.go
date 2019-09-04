package host

import (
	"github.com/journeymidnight/yig-front-caddy"
	"github.com/journeymidnight/yig-front-caddy/caddyhttp/httpserver"
	"strings"
)

func init() {
	caddy.RegisterPlugin("host", caddy.Plugin{
		ServerType: "http",
		Action:     setup,
	})
}

// setup configures a new host middleware instance.
func setup(c *caddy.Controller) error {
	hostConf, err := hostParse(c)
	if err != nil {
		return err
	}
	httpserver.GetConfig(c).AddMiddleware(func(next httpserver.Handler) httpserver.Handler {
		return hostConf
	})
	return nil
}

func hostParse(c *caddy.Controller) (hostConf Host, err error) {
	for c.Next() {
		for c.NextBlock() {
			ext := c.Val()
			switch ext {
			case "s3_db":
				if !c.NextArg() {
					return hostConf, c.ArgErr()
				}
				hostConf.S3Source = c.Val()
				break
			case "business_db":
				if !c.NextArg() {
					return hostConf, c.ArgErr()
				}
				hostConf.BusinessSource = c.Val()
				break
			case "domain":
				if !c.NextArg() {
					return hostConf, c.ArgErr()
				}
				domain := c.Val()
				hostConf.Domain = strings.Split(domain, ",")
				break
			case "custom_domainflag":
				if !c.NextArg() {
					return hostConf, c.ArgErr()
				}
				hostConf.CustomDomainFlag = c.Val()
				break
			case "secret_key":
				if !c.NextArg() {
					return hostConf, c.ArgErr()
				}
				hostConf.SecertKey = c.Val()
			default:
				return hostConf, c.ArgErr()
			}
		}
	}
	return
}
