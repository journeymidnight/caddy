package handle

import "github.com/journeymidnight/yig-front-caddy/caddyhttp/host/meta/types"

type CustomDomainInterface interface {
	//CustomDomain Meta Interface
	GetDomainOfBucketDomain(domainHost string) (info types.DomainInfo, err error)
	ValidBucket(bucket string) (uid string, err error)
	GetDomain(projectId string, domainHost string) (info types.DomainInfo, err error)
	GetDomainInfos(projectId string) (info []types.DomainInfo, err error)
	InsertDomain(customDomainInfo types.DomainInfo) error
	DelDomain(customDomainInfo types.DomainInfo) error
}