package handle

import (
	"encoding/xml"
	"github.com/journeymidnight/yig-front-caddy/caddyhttp/httpserver"
	"github.com/journeymidnight/yig-front-caddy/caddylog"
	"github.com/journeymidnight/yig-front-caddy/helper"
	"net/http"
	"strings"
)

var HOST Host

type Host struct {
	Next             httpserver.Handler
	S3DB             string
	BusinessDB       string
	Domain           []string
	CustomDomainFlag string
	SecretKey        string
	Meta             CustomDomainInterface
	Log              *caddylog.Logger
}

func (h Host) ServeHTTP(w http.ResponseWriter, r *http.Request) (status int, err error) {
	logger := r.Context().Value("logger").(*helper.Log)
	h.Log = logger.Logger
	h.Log.Println(10, r.Method, r.Host, r.Header, r.URL)
	HOST = h
	valid := ValidHost(r.Host)
	v := r.URL.Query()
	flag := v.Get(h.CustomDomainFlag)
	if flag == "" && valid != true {
		status, err := DomainResolution(r)
		if err != nil {
			if status == http.StatusNotFound {
				h.Log.Println(10, status, err)
				response, err := xml.Marshal(GetResponseDomainResolutionErr())
				if err != nil {
					return http.StatusInternalServerError, err
				}
				w.WriteHeader(status)
				return w.Write(response)
			}
			h.Log.Println(10, status, err)
			return status, err
		}
		if status > 300 {
			h.Log.Println(10, status, err)
			return status, err
		}
		h.Log.Println(10, "Custom domain name jump succeeded")
	} else if flag != "" && valid == true {
		claim, status, err := GetMethodFromJWT(r, h.SecretKey)
		if err != nil {
			h.Log.Println(10, status, err)
			return status, err
		}
		Claim = *claim
		result, status, err := DomainOperation(r, flag)
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
	h.Log.Println(10, http.StatusOK, r.Method, r.Host, "Successfully linked yig")
	return h.Next.ServeHTTP(w, r)
}

func ValidHost(host string) bool {
	HOST.Log.Println(10, "Enter ValidHost")
	for _, domain := range HOST.Domain {
		if domain == host {
			return true
		}
		length := len(host) - len(domain)
		if length < 1 {
			continue
		}
		a := strings.LastIndex(host, domain)
		if a == length {
			return true
		}
	}
	return false
}
