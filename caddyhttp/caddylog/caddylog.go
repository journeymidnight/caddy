package caddylog

import (
	"context"
	"github.com/journeymidnight/yig-front-caddy/caddyhttp/httpserver"
	"github.com/journeymidnight/yig-front-caddy/caddylog"
	"github.com/journeymidnight/yig-front-caddy/helper"
	"net/http"
)

var logger *caddylog.Logger

type Log struct {
	Next     httpserver.Handler
	LogPath  string
	LogLevel int
}

func (l Log) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	ctx := context.WithValue(r.Context(), "logger", &helper.Log{Logger: logger})
	return l.Next.ServeHTTP(w, r.WithContext(ctx))
}
