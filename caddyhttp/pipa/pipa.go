package pipa

import (
	"github.com/garyburd/redigo/redis"
	. "github.com/journeymidnight/yig-front-caddy/caddyerrors"
	. "github.com/journeymidnight/yig-front-caddy/caddyhttp/clients/clients/tidbclient"
	"github.com/journeymidnight/yig-front-caddy/caddyhttp/httpserver"
	"github.com/journeymidnight/yig-front-caddy/caddylog"
	"github.com/journeymidnight/yig-front-caddy/helper"
	"net/http"
)

var PIPA Pipa

type Pipa struct {
	Next      httpserver.Handler
	redis     *redis.Pool
	Log       *caddylog.Logger
	SecretKey string
	Client    *TidbClient
}

func (p Pipa) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	logger := r.Context().Value("logger").(*helper.Log)
	p.Log = logger.Logger
	PIPA = p
	key := r.URL.Query().Get("x-oss-process")
	if key != "" {
		p.Log.Println(10, "Enter Pipa Component:", r.Method, r.Host, r.Header, r.URL, "Key:", key)
		var status int
		respose, err := processRequest(r, key)
		if err != nil {
			apiErrorCode, ok := err.(HandleError)
			if ok {
				status = apiErrorCode.HttpStatusCode()
			} else {
				status = http.StatusInternalServerError
			}
			w.WriteHeader(status)
			respose = []byte(apiErrorCode.Description())
			return w.Write(respose)
		}
		w.WriteHeader(http.StatusOK)
		p.Log.Println(10, http.StatusOK, "Picture processed successfully!")
		return w.Write(respose)
	} else {
		p.Log.Println(10, "Pipa:", http.StatusOK, r.Method, r.Host, "Successfully linked yig")
		return p.Next.ServeHTTP(w, r)
	}
}
