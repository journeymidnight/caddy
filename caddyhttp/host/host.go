package host

import (
	"github.com/journeymidnight/yig-front-caddy/caddyhttp/httpserver"
	"net/http"
	"strings"
)

type Host struct {
	Next             httpserver.Handler
	S3Source         string
	BusinessSource   string
	Domain           []string
	CustomDomainFlag string
	SecertKey        string
}

func (h Host) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	// Connect to the database
	conn, err := GetDatabaseConn(h)
	if err != nil {
		return 500, err
	}
	domain := Domain{}
	domain.conn = conn
	host := r.Host
	domain.host = host
	domain.SecertKey = h.SecertKey
	judge := JudgeHost(host, h)
	v := r.URL.Query()
	flag := v.Get(h.CustomDomainFlag)
	if flag == "" && judge != true {
		DomainResolution(w, r, conn)
	} else if flag != "" && judge == true {
		DomainOperation(w, r, domain, flag)
	}
	return h.Next.ServeHTTP(w, r)
}

func JudgeHost(host string, h Host) bool {
	for _, domain := range h.Domain {
		if domain == host {
			return true
		}
		lenth := len(host) - len(domain)
		if lenth < 1 {
			continue
		}
		a := strings.LastIndex(host, domain)
		if a == lenth {
			return true
		}
	}
	return false
}
