package domainresolution

import (
	. "github.com/journeymidnight/yig-front-caddy/caddyerrors"
	"github.com/miekg/dns"
	"net/http"
	"strings"
	"time"
)

var DnsServer = []string{"1.2.4.8", "210.2.4.8", "101.226.4.6", "123.125.81.6", "223.5.5.5", "223.6.6.6", "8.8.8.8", "8.8.4.4", "208.67.222.222", "208.67.220.220", "208.67.222.220", "208.67.220.222"}

func Resolution(r *http.Request) (err error) {
	DOMAINRESOLUTION.Log.Println(10, "Enter domain resolution")
	// Get the corresponding target bucket address
	rHost := r.Host
	validDns, ok := DOMAINRESOLUTION.Cache.Get(rHost)
	if ok != true {
		DOMAINRESOLUTION.Log.Println(20, "Failed to find cache! ")
		var dnsdst []string
		for _, v := range DnsServer {
			dst, err := CNAME(r.Host, v)
			if err != nil {
				return ErrInvalidDnsResolution
			}
			dnsdst = append(dnsdst, dst[0])
		}
		DOMAINRESOLUTION.Log.Println(10, "The domain name resolution address of the CNAME is:", dnsdst)
		for _, h := range dnsdst {
			var host string
			if strings.HasSuffix(h, ".") {
				host = h[0 : len(h)-1]
			} else {
				host = h
			}
			valid, err := DOMAINRESOLUTION.Client.GetDomainOfBucketDomain(rHost)
			if err != nil {
				break
			}
			if valid.DomainBucket == host {
				_ = DOMAINRESOLUTION.Cache.Set(r.Host, host)
				DOMAINRESOLUTION.Log.Println(10, "Insert a key-value pair into the cache:  key =", r.Host, " value =", host)
				r.Host = host
				return nil
			}
		}
		return ErrAccessDenied
	}
	DOMAINRESOLUTION.Log.Println(20, "Succeed to find cache! ")
	DOMAINRESOLUTION.Log.Println(10, "The parameters in the cache are:", validDns)
	domainInfo, err := DOMAINRESOLUTION.Client.GetDomainOfBucketDomain(rHost)
	if err != nil {
		return err
	}
	if domainInfo.DomainBucket != validDns {
		return ErrAccessDenied
	}
	r.Host = domainInfo.DomainBucket
	return nil
}

func CNAME(src string, dnsService string) (dst []string, err error) {
	c := dns.Client{
		Timeout: 5 * time.Second,
	}
	m := dns.Msg{}
	m.SetQuestion(src+".", dns.TypeCNAME)
	r, _, err := c.Exchange(&m, dnsService+":53")
	if err != nil {
		return
	}
	dst = []string{}
	for _, ans := range r.Answer {
		record, isType := ans.(*dns.CNAME)
		if isType {
			dst = append(dst, record.Target)
		}
	}
	return
}
