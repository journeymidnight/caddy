package caddylog

import (
	"github.com/journeymidnight/yig-front-caddy"
	"github.com/journeymidnight/yig-front-caddy/caddyhttp/httpserver"
	"github.com/journeymidnight/yig-front-caddy/caddylog"
	"os"
	"strconv"
)

func init() {
	caddy.RegisterPlugin("caddylog", caddy.Plugin{Action: setup})
}

func setup(c *caddy.Controller) error {
	logPath, logLevel, err := setupCaddyLog(c)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic("Failed to open log file " + logPath)
	}
	logger = caddylog.New(f, "[yig-front-caddy]", caddylog.LstdFlags, logLevel)

	httpserver.GetConfig(c).AddMiddleware(func(next httpserver.Handler) httpserver.Handler {
		return Log{Next: next, LogPath: logPath, LogLevel: logLevel}
	})

	return nil
}

func setupCaddyLog(c *caddy.Controller) (logPath string, level int, err error) {
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
				logLevel := c.Val()
				logLevelInt, err := strconv.Atoi(logLevel)
				if err != nil {
					return logPath, level, err
				}
				level = logLevelInt
				break
			}
		}
	}
	return
}
