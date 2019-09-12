package handle

import (
	"encoding/xml"
	"net/http"
	"time"
)

func DomainOperation(flag string) (response []byte, status int, err error) {
	switch flag {
	case "GetCustomDomain":
		response, status, err = GetCustomDomain()
		return
	case "NewCustomDomain":
		status, err = NewCustomDomain()
		return
	case "DelCustomDomain":
		status, err = DelCustomDomain()
		return
	case "TlsNewCustomDomain":
		status, err = TlsNewCustomDomain()
		return
	case "TlsEditCustomDomain":
		status, err = TlsEditCustomDomain()
		return
	case "TlsDelCustomDomain":
		status, err = TlsDelCustomDomain()
		return
	default:
		status = http.StatusForbidden
		return
	}
}

func GetCustomDomain() ([]byte, int, error) {
	HOST.Log.Logger.Println(10, "Enter get custom domain")
	var data []byte
	projectId := Claim.ProjectId
	object, err := HOST.Meta.GetDomainInfos(projectId)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	var time time.Time
	response := GetResponse(object, time)
	data, err = xml.Marshal(response)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return data, http.StatusOK, nil
}

func NewCustomDomain() (int, error) {
	HOST.Log.Logger.Println(10, "Enter new custom domain")
	projectId := Claim.ProjectId
	domainHost := Claim.DomainHost
	bucket := Claim.Bucket
	domainBucket := Claim.BucketDomain
	uid, err := HOST.Meta.ValidBucket(bucket)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if uid != projectId {
		return http.StatusForbidden, err
	}
	resultHost, err := HOST.Meta.GetDomain(projectId, domainHost)
	if err == nil {
		return http.StatusInternalServerError, err
	}
	if resultHost.DomainHost != "" {
		return http.StatusForbidden, err
	}
	resultHost.ProjectId = projectId
	resultHost.DomainHost = domainHost
	resultHost.DomainBucket = domainBucket
	err = HOST.Meta.InsertDomain(resultHost)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusNoContent, nil
}

func DelCustomDomain() (int, error) {
	HOST.Log.Logger.Println(10, "Enter delete custom domain")
	projectId := Claim.ProjectId
	domainHost := Claim.DomainHost
	info, err := HOST.Meta.GetDomain(projectId, domainHost)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	err = HOST.Meta.DelDomain(info)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusNoContent, nil
}

//TODO
func TlsNewCustomDomain() (status int, err error) {
	status = http.StatusOK
	return
}

//TODO
func TlsEditCustomDomain() (status int, err error) {
	status = http.StatusOK
	return
}

//TODO
func TlsDelCustomDomain() (status int, err error) {
	status = http.StatusOK
	return
}
