package customdomain

import (
	"encoding/xml"
	"github.com/journeymidnight/yig-front-caddy/caddydb/types"
	. "github.com/journeymidnight/yig-front-caddy/caddyerrors"
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
	case "PutCertificate":
		err = PutCertificate(r, claim)
		return
	case "DelCertificate":
		err = DelCertificate(r, claim)
		return
	default:
		return
	}
}

func GetCustomDomain(r *http.Request, claim *Claims) ([]byte, error) {
	if r.Method != "GET" {
		return nil, ErrInvalidRequestMethod
	}
	DOMAIN.Log.Info("Enter get custom domain")
	var data []byte
	projectId := claim.ProjectId
	bucketDomain := claim.BucketDomain
	if projectId == "" {
		return nil, ErrJwtParameterParsing
	}
	var object []types.DomainInfo
	var err error
	if bucketDomain == "" {
		object, err = DOMAIN.CaddyClient.GetAllDomainInfos(projectId, DOMAIN.SealKey)
	} else {
		object, err = DOMAIN.CaddyClient.GetDomainInfos(projectId, bucketDomain, DOMAIN.SealKey)
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
	DOMAIN.Log.Info("Enter new custom domain")
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
	validPID, err := DOMAIN.S3Client.GetBucket(bucket)
	if err != nil {
		return nil, err
	}
	if validPID != projectId {
		return nil, ErrInvalidBucketPermission
	}
	validLength, err := DOMAIN.CaddyClient.GetDomainInfos(projectId, domainBucket, DOMAIN.SealKey)
	if len(validLength) >= 20 {
		return nil, ErrTooManyHostDomainWithBucket
	}
	resultHost, err := DOMAIN.CaddyClient.GetDomain(projectId, domainHost)
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
	err = DOMAIN.CaddyClient.InsertDomain(resultHost, DOMAIN.SealKey)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func DelCustomDomain(r *http.Request, claim *Claims) error {
	if r.Method != "DELETE" {
		return ErrInvalidRequestMethod
	}
	DOMAIN.Log.Info("Enter delete custom domain")
	projectId := claim.ProjectId
	domainHost := claim.DomainHost
	if projectId == "" || domainHost == "" {
		return ErrJwtParameterParsing
	}
	var info types.DomainInfo
	info.ProjectId = projectId
	info.DomainHost = domainHost
	err := DOMAIN.CaddyClient.DelDomain(info)
	if err != nil {
		return err
	}
	return nil
}

func PutCertificate(r *http.Request, claim *Claims) error {
	if r.Method != "POST" {
		return ErrInvalidRequestMethod
	}
	DOMAIN.Log.Info("Enter put custom domain tls")
	projectId := claim.ProjectId
	domainHost := claim.DomainHost
	if projectId == "" || domainHost == "" {
		return ErrJwtParameterParsing
	}
	resultHost, err := DOMAIN.CaddyClient.GetDomain(projectId, domainHost)
	if err != nil {
		return err
	}
	tls, tlsKey, err := GetTlsFromFormData(r)
	if err != nil {
		return err
	}
	err = ValidTls(domainHost, tls, tlsKey)
	if err != nil {
		return err
	}
	resultHost.TlsDomain = tls
	resultHost.TlsDomainKey = tlsKey
	err = DOMAIN.CaddyClient.UpdateDomainTLS(resultHost, DOMAIN.SealKey)
	if err != nil {
		return err
	}
	return nil
}

func DelCertificate(r *http.Request, claim *Claims) (err error) {
	if r.Method != "DELETE" {
		return ErrInvalidRequestMethod
	}
	DOMAIN.Log.Info("Enter delete custom domain tls")
	projectId := claim.ProjectId
	domainHost := claim.DomainHost
	if projectId == "" || domainHost == "" {
		return ErrJwtParameterParsing
	}
	var info types.DomainInfo
	info.ProjectId = projectId
	info.DomainHost = domainHost
	info.TlsDomain = ""
	info.TlsDomainKey = ""
	err = DOMAIN.CaddyClient.UpdateDomainTLS(info, DOMAIN.SealKey)
	if err != nil {
		return err
	}
	return nil
}
