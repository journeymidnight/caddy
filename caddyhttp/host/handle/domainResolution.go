package handle

import (
	"net/http"
)

func DomainResolution(r *http.Request) (status int, err error) {
	HOST.Log.Println(10, "Enter domain resolution")
	// Get the corresponding target bucket address
	domainInfo, err := HOST.Meta.GetDomainOfBucketDomain(r.Host)
	if err != nil {
		return http.StatusForbidden, err
	}
	r.Host = domainInfo.DomainBucket
	return http.StatusOK, nil
}
