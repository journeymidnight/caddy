package customdomain

import (
	"encoding/xml"
	"github.com/journeymidnight/yig-front-caddy/caddyhttp/clients/types"
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
	ProjectId       string `xml:"project_id"`
	HostDomain      string `xml:"host_domain"`
	BucketDomain    string `xml:"bucket_domain"`
	TlsDomainSwitch string `xml:"tls_domain_switch"`
}

type ResponseWithDomainHost struct {
	XMLName            xml.Name `xml:"http://www.unicloud.com CustomDomainResult" json:"-"`
	CustomBucketDomain string
}

func GetResponseWithDomainInfo(data []types.DomainInfo, lastModified time.Time) ResponseWithDomainInfo {
	lastModified = time.Now()
	customDomains := []CustomDomain{}
	for _, data := range data {
		custom := CustomDomain{}
		custom.ProjectId = data.ProjectId
		custom.BucketDomain = data.DomainBucket
		custom.HostDomain = data.DomainHost
		if data.TlsDomain != "" {
			custom.TlsDomainSwitch = "true"
		} else {
			custom.TlsDomainSwitch = "false"
		}
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
