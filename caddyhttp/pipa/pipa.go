package pipa

import (
	"github.com/garyburd/redigo/redis"
	. "github.com/journeymidnight/yig-front-caddy/caddydb/clients/tidbclient"
	. "github.com/journeymidnight/yig-front-caddy/caddyerrors"
	"github.com/journeymidnight/yig-front-caddy/caddyhttp/httpserver"
	"github.com/journeymidnight/yig-front-caddy/caddylog"
	"github.com/journeymidnight/yig-front-caddy/helper"
	"net/http"
	"strings"
)

var PIPA Pipa

var CommonS3ResponseHeaders = []string{"Content-Length", "Content-Type", "Connection", "Date", "ETag", "Server",
	"X-Amz-Delete-Marker", "X-Amz-Request-Id", "X-Amz-Version-Id"}

type Pipa struct {
	Next            httpserver.Handler
	redis           *redis.Pool
	Log             *caddylog.Logger
	SecretKey       string
	S3Client        *TidbClient
	CaddyClient     *TidbClient
	ReservedOrigins []string
}

func (p Pipa) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	logger := r.Context().Value("logger").(*helper.Log)
	p.Log = logger.Logger
	clients := r.Context().Value("database").(map[string]*TidbClient)
	p.S3Client = clients["yig"]
	p.CaddyClient = clients["caddy"]
	PIPA = p
	key := r.URL.Query().Get(IMAGEKEY)
	if key != "" {
		origin := r.Header.Get("Origin")
		if InReservedOrigins(origin) {
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Headers", r.Header.Get("Access-Control-Request-Headers"))
			w.Header().Set("Access-Control-Allow-Methods", r.Header.Get("Access-Control-Request-Method"))
			w.Header().Set("Access-Control-Expose-Headers", strings.Join(CommonS3ResponseHeaders, ","))
		}
		if r.Method == "OPTIONS" {
			if origin == "" || r.Header.Get("Access-Control-Request-Method") == "" {
				w.WriteHeader(ErrInvalidHeader.HttpStatusCode())
				result, _ := writeErrorResponse(ErrInvalidHeader)
				return w.Write(result)
			}
			w.WriteHeader(http.StatusOK)
			return w.Write(nil)
		}
		p.Log.Info("Enter Pipa Component:", r.Method, r.Host, r.Header, r.URL, "Key:", key)
		var status int
		respose, err := imageFunc(r, key)
		if err != nil {
			p.Log.Error(err)
			apiErrorCode, ok := err.(HandleError)
			if ok {
				status = apiErrorCode.HttpStatusCode()
			} else {
				status = http.StatusInternalServerError
			}
			w.WriteHeader(status)
			if respose != nil {
				w.Header().Set("Content-Type", "text/xml;charset=UTF-8")
				return w.Write(respose)
			}
			respose = []byte(apiErrorCode.Description())
			return w.Write(respose)
		}
		w.WriteHeader(http.StatusOK)
		p.Log.Info(http.StatusOK, "Picture processed successfully!")
		return w.Write(respose)
	} else {
		p.Log.Info("Pipa:", http.StatusOK, r.Method, r.Host, "Successfully linked yig")
		return p.Next.ServeHTTP(w, r)
	}
}

func InReservedOrigins(origin string) bool {
	if len(PIPA.ReservedOrigins) == 0 {
		return false
	}
	for _, r := range PIPA.ReservedOrigins {
		if strings.Contains(origin, r) {
			return true
		}
	}
	return false
}
