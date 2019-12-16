package domainresolution

import (
	"github.com/journeymidnight/yig-front-caddy"
	"github.com/journeymidnight/yig-front-caddy/caddyhttp/httpserver"
	"time"
	"zvelo.io/ttlru"
)

func init() {
	caddy.RegisterPlugin("domainresolution", caddy.Plugin{
		ServerType: "http",
		Action:     setup,
	})
}

// setup configures a new host middleware instance.
func setup(c *caddy.Controller) error {
	httpserver.GetConfig(c).AddMiddleware(func(next httpserver.Handler) httpserver.Handler {
		return DomainResolution{
			Next:   next,
			Cache:  ttlru.New(1024, ttlru.WithTTL(10*time.Minute)),
		}
	})
	return nil
}
