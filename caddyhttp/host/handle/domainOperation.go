package handle

import (
	"encoding/xml"
	"fmt"
	"github.com/journeymidnight/yig-front-caddy/caddyhttp/client/types"
	"net/http"
	"strings"
	"time"
)

func DomainOperation(r *http.Request, flag string, claim *Claims) (response []byte, status int, err error) {
	switch flag {
	case "GetCustomDomain":
		response, status, err = GetCustomDomain(r, claim)
		return
	case "NewCustomDomain":
		response, status, err = NewCustomDomain(r, claim)
		return
	case "DelCustomDomain":
		status, err = DelCustomDomain(r, claim)
		return
	case "TlsNewCustomDomain":
		status, err = TlsNewCustomDomain(r, claim)
		return
	case "TlsEditCustomDomain":
		status, err = TlsEditCustomDomain(r, claim)
		return
	case "TlsDelCustomDomain":
		status, err = TlsDelCustomDomain(r, claim)
		return
	default:
		status = http.StatusForbidden
		return
	}
}

func GetCustomDomain(r *http.Request, claim *Claims) ([]byte, int, error) {
	if r.Method != "GET" {
		return nil, http.StatusMethodNotAllowed, fmt.Errorf("The request was made using the wrong request method! ")
	}
	HOST.Log.Println(10, "Enter get custom domain")
	var data []byte
	projectId := claim.ProjectId
	bucketDomain := claim.BucketDomain
	if projectId == "" || bucketDomain == "" {
		return nil, http.StatusBadRequest, fmt.Errorf("Parameter parsing carried by JWT failed. ")
	}
	object, err := HOST.Client.GetDomainInfos(projectId, bucketDomain)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	var time time.Time
	response := GetResponseWithDomainInfo(object, time)
	data, err = xml.Marshal(response)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("Return parameter parsing failed. ")
	}
	return data, http.StatusOK, nil
}

func NewCustomDomain(r *http.Request, claim *Claims) ([]byte, int, error) {
	if r.Method != "PUT" {
		return nil, http.StatusMethodNotAllowed, fmt.Errorf("The request was made using the wrong request method! ")
	}
	HOST.Log.Println(10, "Enter new custom domain")
	projectId := claim.ProjectId
	domainHost := claim.DomainHost
	bucket := claim.Bucket
	domainBucket := claim.BucketDomain
	if projectId == "" || domainHost == "" || bucket == "" || domainBucket == "" {
		return nil, http.StatusBadRequest, fmt.Errorf("Parameter parsing carried by JWT failed. ")
	}
	validDomainBucket := strings.Split(domainBucket, ".")
	if bucket != validDomainBucket[0] {
		return nil, http.StatusBadRequest, fmt.Errorf("Bucket domain name and bucket do not match. ")
	}
	a := strings.HasSuffix(domainBucket, r.Host)
	if a != true {
		return nil, http.StatusPreconditionFailed, fmt.Errorf("The bound domain name does not match the request server domain name. ")
	}
	validPID, err := HOST.Client.GetBucket(bucket)
	if err != nil {
		return nil, http.StatusNotFound, err
	}
	if validPID != projectId {
		return nil, http.StatusNotFound, fmt.Errorf("No bucket operation permission! ")
	}
	resultHost, err := HOST.Client.GetDomain(projectId, domainHost)
	if err != nil {
		HOST.Log.Println(20, "错误是：", err)
	}
	if resultHost.DomainHost != "" {
		response := GetResponseWithDomainHost(resultHost.DomainBucket)
		data, err := xml.Marshal(response)
		if err != nil {
			return nil, http.StatusInternalServerError, fmt.Errorf("Return parameter parsing failed. ")
		}
		return data, http.StatusAlreadyReported, nil
	}
	resultHost.ProjectId = projectId
	resultHost.DomainHost = domainHost
	resultHost.DomainBucket = domainBucket
	err = HOST.Client.InsertDomain(resultHost)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return nil, http.StatusCreated, nil
}

func DelCustomDomain(r *http.Request, claim *Claims) (int, error) {
	if r.Method != "DELETE" {
		return http.StatusMethodNotAllowed, fmt.Errorf("The request was made using the wrong request method! ")
	}
	HOST.Log.Println(10, "Enter delete custom domain")
	projectId := claim.ProjectId
	domainHost := claim.DomainHost
	if projectId == "" || domainHost == "" {
		return http.StatusBadRequest, fmt.Errorf("Parameter parsing carried by JWT failed. ")
	}
	var info types.DomainInfo
	info.ProjectId = projectId
	info.DomainHost = domainHost
	err := HOST.Client.DelDomain(info)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("The specified deleted bucket does not exist! ")
	}
	return http.StatusAccepted, nil
}

//TODO
func TlsNewCustomDomain(r *http.Request, claim *Claims) (status int, err error) {
	status = http.StatusOK
	return
}

//TODO
func TlsEditCustomDomain(r *http.Request, claim *Claims) (status int, err error) {
	status = http.StatusOK
	return
}

//TODO
func TlsDelCustomDomain(r *http.Request, claim *Claims) (status int, err error) {
	status = http.StatusOK
	return
}
