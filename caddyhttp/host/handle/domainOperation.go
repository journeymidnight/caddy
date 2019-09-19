package handle

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func DomainOperation(r *http.Request, flag string) (response []byte, status int, err error) {
	switch flag {
	case "GetCustomDomain":
		response, status, err = GetCustomDomain(r)
		return
	case "NewCustomDomain":
		response, status, err = NewCustomDomain(r)
		return
	case "DelCustomDomain":
		status, err = DelCustomDomain(r)
		return
	case "TlsNewCustomDomain":
		status, err = TlsNewCustomDomain(r)
		return
	case "TlsEditCustomDomain":
		status, err = TlsEditCustomDomain(r)
		return
	case "TlsDelCustomDomain":
		status, err = TlsDelCustomDomain(r)
		return
	default:
		status = http.StatusForbidden
		return
	}
}

func GetCustomDomain(r *http.Request) ([]byte, int, error) {
	if r.Method != "GET" {
		return nil, http.StatusMethodNotAllowed, fmt.Errorf("The request was made using the wrong request method! ")
	}
	HOST.Log.Println(10, "Enter get custom domain")
	var data []byte
	projectId := Claim.ProjectId
	bucketDomain := Claim.BucketDomain
	if projectId == "" || bucketDomain == "" {
		return nil, http.StatusBadRequest, fmt.Errorf("Project_Id cannot be null. ")
	}
	object, err := HOST.Meta.GetDomainInfos(projectId, bucketDomain)
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

func NewCustomDomain(r *http.Request) ([]byte, int, error) {
	if r.Method != "PUT" {
		return nil, http.StatusMethodNotAllowed, fmt.Errorf("The request was made using the wrong request method! ")
	}
	HOST.Log.Println(10, "Enter new custom domain")
	projectId := Claim.ProjectId
	domainHost := Claim.DomainHost
	bucket := Claim.Bucket
	domainBucket := Claim.BucketDomain
	if projectId == "" || domainHost == "" || bucket == "" || domainBucket == "" {
		return nil, http.StatusBadRequest, fmt.Errorf("Parameter parsing carried by JWT failed. ")
	}
	validDomainBucket := strings.Split(domainBucket, ".")
	if bucket != validDomainBucket[0] {
		return nil, http.StatusBadRequest, fmt.Errorf("Bucket domain name and bucket do not match. ")
	}
	length := len(domainBucket) - len(r.Host)
	a := strings.LastIndex(domainBucket, r.Host)
	if a != length {
		return nil, http.StatusPreconditionFailed, fmt.Errorf("The bound domain name does not match the request server domain name. ")
	}
	uid, err := HOST.Meta.ValidBucket(bucket)
	if err != nil {
		return nil, http.StatusNotFound, err
	}
	if uid != projectId {
		return nil, http.StatusNotFound, fmt.Errorf("No bucket operation permission! ")
	}
	resultHost, err := HOST.Meta.GetDomain(projectId, domainHost)
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
	err = HOST.Meta.InsertDomain(resultHost)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return nil, http.StatusCreated, nil
}

func DelCustomDomain(r *http.Request) (int, error) {
	if r.Method != "DELETE" {
		return http.StatusMethodNotAllowed, fmt.Errorf("The request was made using the wrong request method! ")
	}
	HOST.Log.Println(10, "Enter delete custom domain")
	projectId := Claim.ProjectId
	domainHost := Claim.DomainHost
	if projectId == "" || domainHost == "" {
		return http.StatusBadRequest, fmt.Errorf("Parameter parsing carried by JWT failed. ")
	}
	info, err := HOST.Meta.GetDomain(projectId, domainHost)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	err = HOST.Meta.DelDomain(info)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("The specified deleted bucket does not exist! ")
	}
	return http.StatusAccepted, nil
}

//TODO
func TlsNewCustomDomain(r *http.Request) (status int, err error) {
	status = http.StatusOK
	return
}

//TODO
func TlsEditCustomDomain(r *http.Request) (status int, err error) {
	status = http.StatusOK
	return
}

//TODO
func TlsDelCustomDomain(r *http.Request) (status int, err error) {
	status = http.StatusOK
	return
}
