package customdomain

import (
	"encoding/xml"
	. "github.com/journeymidnight/yig-front-caddy/caddyerrors"
	"github.com/journeymidnight/yig-front-caddy/caddyhttp/clients/types"
	"net/http"
	"strings"
	"time"
)

func DomainOperation(r *http.Request, flag string, claim *Claims) (response []byte, err error) {
	switch flag {
	case "GetCustomDomain":
		response, err = GetCustomDomain(r, claim)
		return
	case "NewCustomDomain":
		response, err = NewCustomDomain(r, claim)
		return
	case "DelCustomDomain":
		err = DelCustomDomain(r, claim)
		return
	case "TlsNewCustomDomain":
		err = TlsNewCustomDomain(r, claim)
		return
	case "TlsEditCustomDomain":
		err = TlsEditCustomDomain(r, claim)
		return
	case "TlsDelCustomDomain":
		err = TlsDelCustomDomain(r, claim)
		return
	default:
		return
	}
}

func GetCustomDomain(r *http.Request, claim *Claims) ([]byte, error) {
	if r.Method != "GET" {
		return nil, ErrInvalidRequestMethod
	}
	DOMAIN.Log.Println(10, "Enter get custom domain")
	var data []byte
	projectId := claim.ProjectId
	bucketDomain := claim.BucketDomain
	if projectId == "" {
		return nil, ErrJwtParameterParsing
	}
	var object []types.DomainInfo
	var err error
	if bucketDomain == "" {
		object, err = DOMAIN.Client.GetAllDomainInfos(projectId)
	} else {
		object, err = DOMAIN.Client.GetDomainInfos(projectId, bucketDomain)
	}
	if err != nil {
		return nil, err
	}
	var time time.Time
	response := GetResponseWithDomainInfo(object, time)
	data, err = xml.Marshal(response)
	if err != nil {
		return nil, ErrParameterParsing
	}
	return data, nil
}

func NewCustomDomain(r *http.Request, claim *Claims) ([]byte, error) {
	if r.Method != "PUT" {
		return nil, ErrInvalidRequestMethod
	}
	DOMAIN.Log.Println(10, "Enter new custom domain")
	projectId := claim.ProjectId
	domainHost := claim.DomainHost
	bucket := claim.Bucket
	domainBucket := claim.BucketDomain
	if projectId == "" || domainHost == "" || bucket == "" || domainBucket == "" {
		return nil, ErrJwtParameterParsing
	}
	validDomainBucket := strings.Split(domainBucket, ".")
	if bucket != validDomainBucket[0] {
		return nil, ErrInvalidBucketDomain
	}
	validPID, err := DOMAIN.Client.GetBucket(bucket)
	if err != nil {
		return nil, err
	}
	if validPID != projectId {
		return nil, ErrInvalidBucketPermission
	}
	validLength, err := DOMAIN.Client.GetDomainInfos(projectId, domainBucket)
	if len(validLength) >= 20 {
		return nil, ErrTooManyHostDomainWithBucket
	}
	resultHost, err := DOMAIN.Client.GetDomain(projectId, domainHost)
	if err != nil {
		if err != ErrNoSuchKey {
			return nil, err
		}
	}
	if resultHost.DomainHost != "" {
		response := GetResponseWithDomainHost(resultHost.DomainBucket)
		data, err := xml.Marshal(response)
		if err != nil {
			return nil, ErrGetMarshal
		}
		return data, ErrInvalidHostDomain
	}
	resultHost.ProjectId = projectId
	resultHost.DomainHost = domainHost
	resultHost.DomainBucket = domainBucket
	err = DOMAIN.Client.InsertDomain(resultHost)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func DelCustomDomain(r *http.Request, claim *Claims) error {
	if r.Method != "DELETE" {
		return ErrInvalidRequestMethod
	}
	DOMAIN.Log.Println(10, "Enter delete custom domain")
	projectId := claim.ProjectId
	domainHost := claim.DomainHost
	if projectId == "" || domainHost == "" {
		return ErrJwtParameterParsing
	}
	var info types.DomainInfo
	info.ProjectId = projectId
	info.DomainHost = domainHost
	err := DOMAIN.Client.DelDomain(info)
	if err != nil {
		return err
	}
	return nil
}

//TODO
func TlsNewCustomDomain(r *http.Request, claim *Claims) (err error) {
	return
}

//TODO
func TlsEditCustomDomain(r *http.Request, claim *Claims) (err error) {
	return
}

//TODO
func TlsDelCustomDomain(r *http.Request, claim *Claims) (err error) {
	return
}
