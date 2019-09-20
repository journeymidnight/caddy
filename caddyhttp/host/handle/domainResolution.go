package handle

import (
	"fmt"
	"github.com/miekg/dns"
	"net/http"
	"time"
)

const DNSSERVICE = "114.114.114.114"

func DomainResolution(r *http.Request) (status int, err error) {
	HOST.Log.Println(10, "Enter domain resolution")
	// Get the corresponding target bucket address
	domainInfo, err := HOST.Meta.GetDomainOfBucketDomain(r.Host)
	if err != nil {
		return http.StatusForbidden, fmt.Errorf("No custom domain information was queried! ")
	}
	validDns, ok := HOST.Cache.Get(r.Host)
	if ok != true {
		HOST.Log.Println(20, "Failed to find cache! ")
		dst, err := CNAME(r.Host, DNSSERVICE)
		if err != nil {
			return http.StatusInternalServerError, fmt.Errorf("Query DNS domain name resolution failed. ")
		}
		if len(domainInfo.DomainBucket) > len(dst[0]) {
			HOST.Log.Println(10, "The domain name resolution address of the CNAME is:", dst[0])
			return http.StatusForbidden, fmt.Errorf("No DNS server resolution! ")
		}
		validDns = dst[0][0:len(domainInfo.DomainBucket)]
		HOST.Log.Println(10, "The domain name resolution address of the CNAME is:", validDns)
		if domainInfo.DomainBucket == validDns {
			var test interface{}
			test = validDns
			_ = HOST.Cache.Set(r.Host, test)
			HOST.Log.Println(10, "Insert a key-value pair into the cache:  key =", r.Host, " value =", validDns)
			r.Host = domainInfo.DomainBucket
			return http.StatusOK, nil
		}
	}
	HOST.Log.Println(20, "Succeed to find cache! ")
	HOST.Log.Println(10, "The parameters in the cache are:", validDns)
	if domainInfo.DomainBucket == validDns {
		r.Host = domainInfo.DomainBucket
		return http.StatusOK, nil
	}
	return http.StatusNotFound, fmt.Errorf("No DNS server resolution! ")
}

func CNAME(src string, dnsService string) (dst []string, err error) {
	c := dns.Client{
		Timeout: 5 * time.Second,
	}
	m := dns.Msg{}
	m.SetQuestion(src+".", dns.TypeA)
	r, _, err := c.Exchange(&m, dnsService+":53")
	if err != nil {
		time.Sleep(1 * time.Second * time.Duration(1))
		return
	}
	dst = []string{}
	for _, ans := range r.Answer {
		record, isType := ans.(*dns.CNAME)
		if isType {
			dst = append(dst, record.Target)
		}
	}
	err = nil
	return
}
