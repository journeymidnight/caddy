package handle

import (
	"github.com/journeymidnight/yig-front-caddy/caddyhttp/httpserver"
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
}

func (h Host) ServeHTTP(w http.ResponseWriter, r *http.Request) (status int, err error) {
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
			return w.Write(result)
		}
		return http.StatusOK, nil
	}
	return h.Next.ServeHTTP(w, r)
}

func ValidHost(host string) bool {
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
