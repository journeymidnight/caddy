package pipa

import (
	"github.com/garyburd/redigo/redis"
	"github.com/journeymidnight/yig-front-caddy/caddyerrors"
	"github.com/journeymidnight/yig-front-caddy/caddyhttp/httpserver"
	"github.com/journeymidnight/yig-front-caddy/caddylog"
	"github.com/journeymidnight/yig-front-caddy/helper"
	"net/http"
)

var PIPA Pipa

type Pipa struct {
	Next  httpserver.Handler
	redis *redis.Pool
	Log   *caddylog.Logger
}

func (p Pipa) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	logger := r.Context().Value("logger").(*helper.Log)
	p.Log = logger.Logger
	PIPA = p
	key := r.URL.Query().Get("x-oss-process")
	if key != "" {
		p.Log.Println(10, "Enter Pipa Component:", r.Method, r.Host, r.Header, r.URL)
		respose, err := processRequest(w, r, key)
		if err != nil {
			switch err {
			case caddyerrors.ErrNoRouter:
				w.WriteHeader(http.StatusNotFound)
			case caddyerrors.ErrTimeout:
				w.WriteHeader(http.StatusRequestTimeout)
			case caddyerrors.ErrInternalServer:
				w.WriteHeader(http.StatusInternalServerError)
			default:
				w.WriteHeader(http.StatusNotFound)
			}
			return w.Write(nil)
		}
		w.WriteHeader(http.StatusOK)
		p.Log.Println(10, http.StatusOK, "Picture processed successfully!")
		return w.Write(respose)
	} else {
		p.Log.Println(10, "Pipa:", http.StatusOK, r.Method, r.Host, "Successfully linked yig")
		return p.Next.ServeHTTP(w, r)
	}
}
