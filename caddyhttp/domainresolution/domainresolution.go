package domainresolution

import (
	. "github.com/journeymidnight/yig-front-caddy/caddyerrors"
	"github.com/journeymidnight/yig-front-caddy/caddyhttp/clients/clients/tidbclient"
	"github.com/journeymidnight/yig-front-caddy/caddyhttp/httpserver"
	"github.com/journeymidnight/yig-front-caddy/caddylog"
	"github.com/journeymidnight/yig-front-caddy/helper"
	"net/http"
	"zvelo.io/ttlru"
)

var DOMAINRESOLUTION DomainResolution

type DomainResolution struct {
	Next         httpserver.Handler
	Client       *tidbclient.TidbClient
	Log          *caddylog.Logger
	Cache        ttlru.Cache
	TlsSecretKey string
}

func (c DomainResolution) ServeHTTP(w http.ResponseWriter, r *http.Request) (status int, err error) {
	logger := r.Context().Value("logger").(*helper.Log)
	c.Log = logger.Logger
	c.Log.Println(10, "Enter DomainResolution Function", r.Method, r.Host, r.Header, r.URL)
	DOMAINRESOLUTION = c
	err = Resolution(r)
	if err != nil {
		c.Log.Println(10, err)
		apiErrorCode, ok := err.(HandleError)
		if ok {
			status = apiErrorCode.HttpStatusCode()
		} else {
			status = http.StatusInternalServerError
		}
		c.Log.Println(10, status, err)
		return status, err
	}
	c.Log.Println(10, http.StatusOK, r.Method, r.Host, "Successfully linked yig")
	return c.Next.ServeHTTP(w, r)
}
