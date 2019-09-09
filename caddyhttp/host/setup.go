package host

import (
	"github.com/journeymidnight/yig-front-caddy"
	"github.com/journeymidnight/yig-front-caddy/caddyhttp/host/handle"
	"github.com/journeymidnight/yig-front-caddy/caddyhttp/host/meta"
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
	s3DB, businessDB, domain, customDomainFlag, secretKey, err := hostParse(c)
	if err != nil {
		return err
	}
	httpserver.GetConfig(c).AddMiddleware(func(next httpserver.Handler) httpserver.Handler {
		return handle.Host{
			Next:             next,
			S3DB:             s3DB,
			BusinessDB:       businessDB,
			Domain:           domain,
			CustomDomainFlag: customDomainFlag,
			SecretKey:        secretKey,
			Meta:             meta.New(s3DB, businessDB),
		}
	})
	return nil
}

func hostParse(c *caddy.Controller) (s3DB string, businessDB string, domainInfo []string, customDomainFlag string, secretKey string, err error) {
	for c.Next() {
		for c.NextBlock() {
			ext := c.Val()
			switch ext {
			case "s3_db":
				if !c.NextArg() {
					return
				}
				s3DB = c.Val()
				break
			case "business_db":
				if !c.NextArg() {
					return
				}
				businessDB = c.Val()
				break
			case "domain":
				if !c.NextArg() {
					return
				}
				domain := c.Val()
				domainInfo = strings.Split(domain, ",")
				break
			case "custom_domainflag":
				if !c.NextArg() {
					return
				}
				customDomainFlag = c.Val()
				break
			case "secret_key":
				if !c.NextArg() {
					return
				}
				secretKey = c.Val()
			default:
				return
			}
		}
	}
	return
}
