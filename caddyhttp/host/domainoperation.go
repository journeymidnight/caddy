package host

import (
	"encoding/xml"
	"net/http"
	"strconv"
	"time"
)

type Domain struct {
	conn       *TidbClient
	host       string
	customInfo CustomInfo
	tlsInfo    TlsInfo
	editBucket string
	SecertKey  string
}

type CustomInfo struct {
	projectId  string
	bucket     string
	domain     string
	domainHost string
}

type TlsInfo struct {
	tls    string
	tlsKey string
}

type CustomDomainInfo struct {
	CustomDomain []CustomDomain
}

type CustomDomain struct {
	ProjectId  string `xml:"project-id"`
	DomainHost string `xml:"host"`
	Domain     string `xml:"bucket-domain"`
	TlsKey     string `xml:"tls-key"`
	Tls        string `xml:"tls"`
}

func DomainOperation(w http.ResponseWriter, r *http.Request, domain Domain, flag string) {
	claim := GetMethodFromJWT(w, r, domain.SecertKey)
	bucket_domain := claim.Bucket + "." + domain.host
	domain.customInfo.projectId = claim.ProjectId
	domain.customInfo.bucket = claim.Bucket
	domain.customInfo.domainHost = claim.DomainHost
	domain.customInfo.domain = bucket_domain
	switch flag {
	case "getcustomdomain":
		GetCustomDomain(w, r, domain)
	case "newcustomdomain":
		NewCustomDomain(w, r, domain)
	case "editcustomdomain":
		editbucket := r.Header.Get("edit-bucket")
		domain.editBucket = editbucket
		EditCustomDomain(w, r, domain)
	case "delcustomdomain":
		DelCustomDomain(w, r, domain)
	case "tlsnewcustomdomain":
		TlsNewCustomDomain(w, r, domain)
	case "tlseditcustomdomain":
		TlsEditCustomDomain(w, r, domain)
	case "tlsdelcustomdomain":
		TlsdelCustomDomain(w, r, domain)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func GetCustomDomain(w http.ResponseWriter, r *http.Request, domain Domain) {
	var data, result []byte
	object, err := GetDomainInfo(domain)
	if err != nil {
		w.WriteHeader(http.StatusRequestTimeout)
	}
	data, err = xml.Marshal(object)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	var time time.Time
	w.WriteHeader(http.StatusOK)
	response := GenerateObjectResponse(time)
	encodedSuccessResponse := EncodeResponse(response)
	result = append(result, encodedSuccessResponse...)
	result = append(result, data...)
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.Write(result)
}

func NewCustomDomain(w http.ResponseWriter, r *http.Request, domain Domain) {
	resultuser, err := AuthBucketInfo(domain)
	if err != nil || resultuser[0] == "" {
		w.WriteHeader(http.StatusForbidden)
	}
	resulthost, err := AuthHostInfo(domain)
	if err != nil || resulthost != nil {
		w.WriteHeader(http.StatusForbidden)
	}
	err = InsertDomainInfo(domain)
	if err != nil {
		w.WriteHeader(http.StatusRequestTimeout)
	}
	var time time.Time
	w.WriteHeader(http.StatusOK)
	response := GenerateObjectResponse(time)
	encodedSuccessResponse := EncodeResponse(response)
	w.Write(encodedSuccessResponse)
}

func EditCustomDomain(w http.ResponseWriter, r *http.Request, domain Domain) {
	resultuser, err := AuthBucketInfo(domain)
	if err != nil || resultuser[0] == "" {
		w.WriteHeader(http.StatusForbidden)
	}
	rows, err := AuthHostInfo(domain)
	if err != nil || rows == nil {
		w.WriteHeader(http.StatusForbidden)
	}
	if domain.editBucket != "" {
		editBucket := domain.editBucket + "." + domain.host
		err = EditDomainInfo(domain, editBucket)
		if err != nil {
			w.WriteHeader(http.StatusRequestTimeout)
		}
	} else {
		w.WriteHeader(http.StatusForbidden)
	}
	var time time.Time
	w.WriteHeader(http.StatusOK)
	response := GenerateObjectResponse(time)
	encodedSuccessResponse := EncodeResponse(response)
	w.Write(encodedSuccessResponse)
}

func DelCustomDomain(w http.ResponseWriter, r *http.Request, domain Domain) {
	rows, err := AuthHostInfo(domain)
	if err != nil || rows == nil {
		w.WriteHeader(http.StatusForbidden)
	}
	err = DelDomainInfo(domain)
	if err != nil {
		w.WriteHeader(http.StatusRequestTimeout)
	}
	var time time.Time
	w.WriteHeader(http.StatusOK)
	response := GenerateObjectResponse(time)
	encodedSuccessResponse := EncodeResponse(response)
	w.Write(encodedSuccessResponse)
}

//TODO
func TlsNewCustomDomain(w http.ResponseWriter, r *http.Request, domain Domain) {
	var time time.Time
	w.WriteHeader(http.StatusOK)
	response := GenerateObjectResponse(time)
	encodedSuccessResponse := EncodeResponse(response)
	w.Write(encodedSuccessResponse)
}

//TODO
func TlsEditCustomDomain(w http.ResponseWriter, r *http.Request, domain Domain) {
	var time time.Time
	w.WriteHeader(http.StatusOK)
	response := GenerateObjectResponse(time)
	encodedSuccessResponse := EncodeResponse(response)
	w.Write(encodedSuccessResponse)
}

//TODO
func TlsdelCustomDomain(w http.ResponseWriter, r *http.Request, domain Domain) {
	var time time.Time
	w.WriteHeader(http.StatusOK)
	response := GenerateObjectResponse(time)
	encodedSuccessResponse := EncodeResponse(response)
	w.Write(encodedSuccessResponse)
}
