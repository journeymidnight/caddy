package handle

import (
	"github.com/journeymidnight/yig-front-caddy/caddyhttp/client/clients/tidbclient"
	"github.com/journeymidnight/yig-front-caddy/caddyhttp/httpserver"
	"github.com/journeymidnight/yig-front-caddy/caddylog"
	"github.com/journeymidnight/yig-front-caddy/helper"
	"net/http"
	"zvelo.io/ttlru"
)

var HOST Host

type Host struct {
	Next             httpserver.Handler
	CustomDomainFlag string
	SecretKey        string
	Client           *tidbclient.TidbClient
	Log              *caddylog.Logger
	Cache            ttlru.Cache
}

func (h Host) ServeHTTP(w http.ResponseWriter, r *http.Request) (status int, err error) {
	logger := r.Context().Value("logger").(*helper.Log)
	h.Log = logger.Logger
	h.Log.Println(10, r.Method, r.Host, r.Header, r.URL)
	HOST = h
	v := r.URL.Query()
	flag := v.Get(h.CustomDomainFlag)
	if flag == "" {
		status, err := DomainResolution(r)
		if err != nil {
			h.Log.Println(10, status, err)
			return status, err
		}
		if status > 300 {
			h.Log.Println(10, status, err)
			return status, err
		}
		h.Log.Println(10, http.StatusOK, r.Method, r.Host, "Successfully linked yig")
		return h.Next.ServeHTTP(w, r)
	}
	var claim *Claims
	claim, status, err = GetMethodFromJWT(r, h.SecretKey)
	if err != nil {
		h.Log.Println(10, status, err)
		return status, err
	}
	result, status, err := DomainOperation(r, flag, claim)
	if err != nil || status > 300 {
		h.Log.Println(10, status, err)
		return status, err
	}
	if result != nil {
		w.WriteHeader(status)
		h.Log.Println(10, status, "The information returned is:", string(result))
		return w.Write(result)
	}
	h.Log.Println(10, status, "Custom domain name succeeded")
	w.WriteHeader(status)
	return w.Write(result)
}
