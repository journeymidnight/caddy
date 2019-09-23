package host

import (
	"github.com/journeymidnight/yig-front-caddy"
	"github.com/journeymidnight/yig-front-caddy/caddyhttp/client/clients/tidbclient"
	"github.com/journeymidnight/yig-front-caddy/caddyhttp/host/handle"
	"github.com/journeymidnight/yig-front-caddy/caddyhttp/httpserver"
	"strconv"
	"time"
	"zvelo.io/ttlru"
)

func init() {
	caddy.RegisterPlugin("host", caddy.Plugin{
		ServerType: "http",
		Action:     setup,
	})
}

// setup configures a new host middleware instance.
func setup(c *caddy.Controller) error {
	db, customDomainFlag, secretKey, err := hostParse(c)
	if err != nil {
		return err
	}
	httpserver.GetConfig(c).AddMiddleware(func(next httpserver.Handler) httpserver.Handler {
		return handle.Host{
			Next:             next,
			CustomDomainFlag: customDomainFlag,
			SecretKey:        secretKey,
			Client:           tidbclient.NewTidbClient(db),
			Cache:            ttlru.New(1024, ttlru.WithTTL(10*time.Minute)),
		}
	})
	return nil
}

func hostParse(c *caddy.Controller) (db tidbclient.DBInfo, customDomainFlag string, secretKey string, err error) {
	for c.Next() {
		for c.NextBlock() {
			ext := c.Val()
			switch ext {
			case "s3_db":
				if !c.NextArg() {
					return db, "", "", c.ArgErr()
				}
				db.S3Source = c.Val()
				break
			case "caddy_db":
				if !c.NextArg() {
					return db, "", "", c.ArgErr()
				}
				db.CaddySource = c.Val()
				break
			case "db_max_idle_conns":
				if !c.NextArg() {
					return db, "", "", c.ArgErr()
				}
				int, err := strconv.Atoi(c.Val())
				if err != nil {
					return db, "", "", err
				}
				db.DBMaxIdleConns = int
			case "db_max_open_conns":
				if !c.NextArg() {
					return db, "", "", c.ArgErr()
				}
				int, err := strconv.Atoi(c.Val())
				if err != nil {
					return db, "", "", err
				}
				db.DBMaxOpenConns = int
			case "db_conn_max_life_seconds":
				if !c.NextArg() {
					return db, "", "", c.ArgErr()
				}
				int, err := strconv.Atoi(c.Val())
				if err != nil {
					return db, "", "", err
				}
				db.DBConnMaxLifeSeconds = int
			case "custom_domainflag":
				if !c.NextArg() {
					return db, "", "", c.ArgErr()
				}
				customDomainFlag = c.Val()
				break
			case "secret_key":
				if !c.NextArg() {
					return db, "", "", c.ArgErr()
				}
				secretKey = c.Val()
			default:
				return db, "", "", c.ArgErr()
			}
		}
	}
	return
}
