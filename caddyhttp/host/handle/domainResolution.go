package handle

import (
	"net/http"
)

func DomainResolution(r *http.Request) (status int, err error) {
	HOST.Log.Logger.Println(10, "Enter domain resolution")
	// Get the target bucket
	projectID := r.URL.Query().Get("project_id")
	if projectID == "" {
		return http.StatusForbidden, err
	}
	// Get the corresponding target bucket address
	domainInfo, err := HOST.Meta.GetDomainOfBucketDomain(projectID, r.Host)
	if err != nil {
		return http.StatusForbidden, err
	}
	r.Host = domainInfo.DomainBucket
	return http.StatusOK, nil
}
