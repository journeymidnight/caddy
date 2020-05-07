package caddylog

import (
	"github.com/journeymidnight/yig-front-caddy"
	"github.com/journeymidnight/yig-front-caddy/caddyhttp/httpserver"
	log "github.com/journeymidnight/yig-front-caddy/caddylog"
)

func init() {
	caddy.RegisterPlugin("caddylog", caddy.Plugin{Action: setup})
}

func setup(c *caddy.Controller) error {
	logPath, logLevel, err := setupCaddyLog(c)
	if err != nil {
		return err
	}
	level := log.ParseLevel(logLevel)
	logger = log.NewFileLogger(logPath, level)
	httpserver.GetConfig(c).AddMiddleware(func(next httpserver.Handler) httpserver.Handler {
		return Log{Next: next, LogPath: logPath, LogLevel: logLevel}
	})

	return nil
}

func setupCaddyLog(c *caddy.Controller) (logPath string, level string, err error) {
	for c.Next() {
		for c.NextBlock() {
			switch c.Val() {
			case "log_path":
				if !c.NextArg() {
					return logPath, level, c.ArgErr()
				}
				logPath = c.Val()
				break
			case "log_level":
				if !c.NextArg() {
					return logPath, level, c.ArgErr()
				}
				level = c.Val()
				break
			}
		}
	}
	return
}
