package customdomain

import (
	"github.com/journeymidnight/yig-front-caddy/caddydb/clients/tidbclient"
	. "github.com/journeymidnight/yig-front-caddy/caddyerrors"
	"github.com/journeymidnight/yig-front-caddy/caddyhttp/httpserver"
	"github.com/journeymidnight/yig-front-caddy/caddylog"
	"github.com/journeymidnight/yig-front-caddy/helper"
	"net/http"
	"zvelo.io/ttlru"
)

var DOMAIN Domain

type Domain struct {
	Next             httpserver.Handler
	CustomDomainFlag string
	SecretKey        string
	SealKey          string
	S3Client         *tidbclient.TidbClient
	CaddyClient      *tidbclient.TidbClient
	Log              *caddylog.Logger
	Cache            ttlru.Cache
}

func (c Domain) ServeHTTP(w http.ResponseWriter, r *http.Request) (status int, err error) {
	logger := r.Context().Value("logger").(*helper.Log)
	c.Log = logger.Logger
	clients := r.Context().Value("database").(map[string]*tidbclient.TidbClient)
	c.S3Client = clients["yig"]
	c.CaddyClient = clients["caddy"]
	DOMAIN = c
	v := r.URL.Query()
	flag := v.Get(c.CustomDomainFlag)
	if flag == "" {
		c.Log.Info("Customdomain:", http.StatusOK, r.Method, r.Host, "Successfully linked yig")
		return c.Next.ServeHTTP(w, r)
	}
	c.Log.Info("Enter CustomDomain Function", r.Method, r.Host, r.Header, r.URL)
	var claim *Claims
	claim, err = GetMethodFromJWT(r, c.SecretKey)
	if err != nil {
		apiErrorCode, ok := err.(HandleError)
		if ok {
			status = apiErrorCode.HttpStatusCode()
		} else {
			status = http.StatusInternalServerError
		}
		c.Log.Error(status, err)
		return status, err
	}
	result, err := DomainOperation(r, flag, claim)
	if err != nil {
		c.Log.Println(10, err)
		apiErrorCode, ok := err.(HandleError)
		if ok {
			status = apiErrorCode.HttpStatusCode()
		} else {
			status = http.StatusInternalServerError
		}
		if err != ErrInvalidHostDomain {
			c.Log.Error(status, err)
			return status, err
		} else {
			w.WriteHeader(status)
			c.Log.Error(status, "The information returned is:", string(result))
			return w.Write(result)
		}
	}
	if result != nil {
		w.WriteHeader(http.StatusOK)
		c.Log.Info(http.StatusOK, "The information returned is:", string(result))
		return w.Write(result)
	}
	c.Log.Info(http.StatusOK, "Custom domain name succeeded")
	w.WriteHeader(http.StatusOK)
	return w.Write(result)
}
