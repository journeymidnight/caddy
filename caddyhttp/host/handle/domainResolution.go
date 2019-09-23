package handle

import (
	"fmt"
	"github.com/miekg/dns"
	"net/http"
	"strings"
	"time"
)

const DNSSERVICE = "114.114.114.114"

func DomainResolution(r *http.Request) (status int, err error) {
	HOST.Log.Println(10, "Enter domain resolution")
	// Get the corresponding target bucket address
	rHost := r.Host
	validDns, ok := HOST.Cache.Get(rHost)
	if ok != true {
		HOST.Log.Println(20, "Failed to find cache! ")
		dst, err := CNAME(r.Host, DNSSERVICE)
		if err != nil {
			return http.StatusInternalServerError, fmt.Errorf("Query DNS domain name resolution failed. ")
		}
		HOST.Log.Println(10, "The domain name resolution address of the CNAME is:", dst)
		for _, h := range dst {
			var host string
			if strings.HasSuffix(h, ".") {
				host = h[0 : len(h)-1]
			} else {
				host = h
			}
			valid, err := HOST.Client.GetDomainOfBucketDomain(rHost)
			if err != nil {
				break
			}
			if valid.DomainBucket == host {
				_ = HOST.Cache.Set(r.Host, host)
				HOST.Log.Println(10, "Insert a key-value pair into the cache:  key =", r.Host, " value =", host)
				r.Host = host
				return http.StatusOK, nil
			}
		}
		return http.StatusNotFound, fmt.Errorf("No DNS server resolution! ")
	}
	HOST.Log.Println(20, "Succeed to find cache! ")
	HOST.Log.Println(10, "The parameters in the cache are:", validDns)
	domainInfo, err := HOST.Client.GetDomainOfBucketDomain(rHost)
	if err != nil {
		return http.StatusForbidden, fmt.Errorf("No custom domain information was queried! ")
	}
	if domainInfo.DomainBucket != validDns {
		return http.StatusNotFound, fmt.Errorf("No DNS server resolution! ")
	}
	r.Host = domainInfo.DomainBucket
	return http.StatusOK, nil
}

func CNAME(src string, dnsService string) (dst []string, err error) {
	c := dns.Client{
		Timeout: 5 * time.Second,
	}
	m := dns.Msg{}
	m.SetQuestion(src+".", dns.TypeA)
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
