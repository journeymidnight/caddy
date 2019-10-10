package domainresolution

import (
	"github.com/journeymidnight/yig-front-caddy"
	"github.com/journeymidnight/yig-front-caddy/caddyhttp/clients/clients/tidbclient"
	"github.com/journeymidnight/yig-front-caddy/caddyhttp/httpserver"
	"strconv"
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
	s3Source, caddySource, db, tlsSecretKey, err := customParse(c)
	if err != nil {
		return err
	}
	httpserver.GetConfig(c).AddMiddleware(func(next httpserver.Handler) httpserver.Handler {
		return DomainResolution{
			Next:         next,
			Client:       tidbclient.NewCustomDomainClient(s3Source, caddySource, db),
			Cache:        ttlru.New(1024, ttlru.WithTTL(10*time.Minute)),
			TlsSecretKey: tlsSecretKey,
		}
	})
	return nil
}

func customParse(c *caddy.Controller) (s3Source string, caddySource string, db tidbclient.DBInfo, tlsSecretKey string, err error) {
	for c.Next() {
		for c.NextBlock() {
			ext := c.Val()
			switch ext {
			case "s3_db":
				if !c.NextArg() {
					return s3Source, "", db, "", c.ArgErr()
				}
				s3Source = c.Val()
				break
			case "caddy_db":
				if !c.NextArg() {
					return "", caddySource, db, "", c.ArgErr()
				}
				caddySource = c.Val()
				break
			case "db_max_idle_conns":
				if !c.NextArg() {
					return "", "", db, "", c.ArgErr()
				}
				int, err := strconv.Atoi(c.Val())
				if err != nil {
					return "", "", db, "", err
				}
				db.DBMaxIdleConns = int
			case "db_max_open_conns":
				if !c.NextArg() {
					return "", "", db, "", c.ArgErr()
				}
				int, err := strconv.Atoi(c.Val())
				if err != nil {
					return "", "", db, "", err
				}
				db.DBMaxOpenConns = int
			case "db_conn_max_life_seconds":
				if !c.NextArg() {
					return "", "", db, "", c.ArgErr()
				}
				int, err := strconv.Atoi(c.Val())
				if err != nil {
					return "", "", db, "", err
				}
				db.DBConnMaxLifeSeconds = int
			case "tls_secret_key":
				if !c.NextArg() {
					return "", "", db, "", c.ArgErr()
				}
				tlsSecretKey = c.Val()
			default:
				return "", "", db, "", c.ArgErr()
			}
		}
	}
	return
}
