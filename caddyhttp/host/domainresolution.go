package host

import (
	"net/http"
)

func DomainResolution(w http.ResponseWriter, r *http.Request, conn *TidbClient) {
	// Get the target bucket
	host := r.Host
	projectID := r.Header.Get("project_id")
	// Get the corresponding target bucket address
	bucket_domain, err := GetTargetBucket(conn, host, projectID)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
	}
	r.Host = bucket_domain
}
