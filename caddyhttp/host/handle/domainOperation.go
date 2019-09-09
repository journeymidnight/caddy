package handle

import (
	"encoding/xml"
	"net/http"
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
		status, err = TlsdelCustomDomain()
		return
	default:
		status = http.StatusForbidden
		return
	}
}

func GetCustomDomain() ([]byte, int, error) {
	var data []byte
	projectId := Claim.ProjectId
	object, err := HOST.Meta.GetDomainInfos(projectId)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	data, err = xml.Marshal(object)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return data, http.StatusOK, nil
}

func NewCustomDomain() (int, error) {
	projectId := Claim.ProjectId
	domainhost := Claim.DomainHost
	bucket := Claim.Bucket
	domainbucket := Claim.BucketDomain
	uid, err := HOST.Meta.ValidBucket(bucket)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if uid != projectId {
		return http.StatusForbidden, err
	}
	resulthost, err := HOST.Meta.GetDomain(projectId, domainhost)
	if err == nil {
		return http.StatusInternalServerError, err
	}
	if resulthost.DomainHost != "" {
		return http.StatusForbidden, err
	}
	resulthost.ProjectId = projectId
	resulthost.DomainHost = domainhost
	resulthost.DomainBucket = domainbucket
	err = HOST.Meta.InsertDomain(resulthost)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusNoContent, nil
}

func DelCustomDomain() (int, error) {
	projectId := Claim.ProjectId
	domainhost := Claim.DomainHost
	info, err := HOST.Meta.GetDomain(projectId, domainhost)
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
func TlsdelCustomDomain() (status int, err error) {
	status = http.StatusOK
	return
}
