package handle

import (
	"github.com/journeymidnight/yig-front-caddy/caddyhttp/httpserver"
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
	Log              *helper.Log
}

func (h Host) ServeHTTP(w http.ResponseWriter, r *http.Request) (status int, err error) {
	logger := r.Context().Value("logger").(*helper.Log)
	h.Log = logger
	h.Log.Logger.Println(10, r.Host, r.Header, r.URL)
	HOST = h
	valid := ValidHost(r.Host)
	v := r.URL.Query()
	flag := v.Get(h.CustomDomainFlag)
	if flag == "" && valid != true {
		status, err := DomainResolution(r)
		if err != nil {
			return status, err
		}
		if status != http.StatusOK {
			return status, err
		}
		h.Log.Logger.Println(10, http.StatusOK, "Custom domain name jump succeeded")
	} else if flag != "" && valid == true {
		claim, status, err := GetMethodFromJWT(r, h.SecretKey)
		if err != nil {
			return status, err
		}
		Claim = *claim
		result, status, err := DomainOperation(flag)
		if err != nil {
			return status, err
		}
		if result != nil {
			w.WriteHeader(http.StatusOK)
			h.Log.Logger.Println(10, http.StatusOK, "Get custom domain success:", result)
			return w.Write(result)
		}
		h.Log.Logger.Println(10, http.StatusOK, "Custom domain name succeeded")
		return http.StatusOK, nil
	}
	h.Log.Logger.Println(10, http.StatusOK, "Successfully linked yig")
	return h.Next.ServeHTTP(w, r)
}

func ValidHost(host string) bool {
	HOST.Log.Logger.Println(10, "Enter ValidHost")
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
