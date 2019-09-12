package caddylog

import (
	"context"
	"github.com/journeymidnight/yig-front-caddy/caddyhttp/httpserver"
	"github.com/journeymidnight/yig-front-caddy/caddylog"
	"github.com/journeymidnight/yig-front-caddy/helper"
	"net/http"
	"os"
)

var logger *caddylog.Logger

type Log struct {
	Next     httpserver.Handler
	LogPath  string
	LogLevel int
}

func (l Log) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	f, err := os.OpenFile(l.LogPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic("Failed to open log file " + l.LogPath)
	}

	logger = caddylog.New(f, "[yig-front-caddy]", caddylog.LstdFlags, l.LogLevel)
	helper.Logger = logger
	ctx := context.WithValue(r.Context(), "logger", &helper.Log{Logger: logger})
	return l.Next.ServeHTTP(w, r.WithContext(ctx))
}
