package domainresolution

import (
	. "github.com/journeymidnight/yig-front-caddy/caddyerrors"
	"github.com/miekg/dns"
	"net/http"
	"strings"
	"time"
)

const DnsServer = "114.114.114.114"

func Resolution(r *http.Request) (err error) {
	DOMAINRESOLUTION.Log.Info("Enter domain resolution")
	// Get the corresponding target bucket address
	rHost := r.Host
	validDns, ok := DOMAINRESOLUTION.Cache.Get(rHost)
	if ok != true {
		DOMAINRESOLUTION.Log.Info("Failed to find cache! ")
		dst, err := CNAME(r.Host, DnsServer)
		if err != nil {
			return ErrInvalidDnsResolution
		}
		DOMAINRESOLUTION.Log.Info("The domain name resolution address of the CNAME is:", dst)
		for _, h := range dst {
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
				DOMAINRESOLUTION.Log.Info("Insert a key-value pair into the cache:  key =", r.Host, " value =", host)
				r.Host = host
				return nil
			}
		}
		return ErrAccessDenied
	}
	DOMAINRESOLUTION.Log.Info("Succeed to find cache! ")
	DOMAINRESOLUTION.Log.Info("The parameters in the cache are:", validDns)
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
