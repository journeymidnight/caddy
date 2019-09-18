package handle

import (
	"encoding/xml"
	"github.com/journeymidnight/yig-front-caddy/caddyhttp/host/meta/types"
	"time"
)

const (
	timeFormatAMZ = "2006-01-02T15:04:05.000Z" // Reply date format
)

type ResponseWithDomainInfo struct {
	XMLName          xml.Name `xml:"http://www.unicloud.com CustomDomainResult" json:"-"`
	LastModified     string   // time string of format "2006-01-02T15:04:05.000Z"
	CustomDomainInfo CustomDomainInfo
}

type CustomDomainInfo struct {
	CustomDomain []CustomDomain
}

type CustomDomain struct {
	ProjectId    string `xml:"project_id"`
	HostDomain   string `xml:"host_domain"`
	BucketDomain string `xml:"bucket_domain"`
}

type ResponseWithDomainHost struct {
	XMLName            xml.Name `xml:"http://www.unicloud.com CustomDomainResult" json:"-"`
	CustomBucketDomain string
}

type ResponseDomainResolutionErr struct {
	XMLName xml.Name `xml:"http://www.unicloud.com CustomDomainResult" json:"-"`
	Reason  string
}

func GetResponseWithDomainInfo(data []types.DomainInfo, lastModified time.Time) ResponseWithDomainInfo {
	lastModified = time.Now()
	customDomains := []CustomDomain{}
	for _, data := range data {
		custom := CustomDomain{}
		custom.ProjectId = data.ProjectId
		custom.BucketDomain = data.DomainBucket
		custom.HostDomain = data.DomainHost
		customDomains = append(customDomains, custom)
	}
	customDomainInfo := CustomDomainInfo{}
	customDomainInfo.CustomDomain = customDomains
	return ResponseWithDomainInfo{
		LastModified:     lastModified.UTC().Format(timeFormatAMZ),
		CustomDomainInfo: customDomainInfo,
	}
}

func GetResponseWithDomainHost(data string) ResponseWithDomainHost {
	return ResponseWithDomainHost{
		CustomBucketDomain: data,
	}
}

func GetResponseDomainResolutionErr() ResponseDomainResolutionErr {
	return ResponseDomainResolutionErr{
		Reason: "The domain name you are using does not have the corresponding CNAME domain name resolution or the domain name CNAME resolution has not been bound.",
	}
}