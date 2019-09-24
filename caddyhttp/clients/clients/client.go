package clients

import (
	"github.com/journeymidnight/yig-front-caddy/caddyhttp/clients/types"
)

type Client interface {
	//Transaction
	NewTrans() (tx interface{}, err error)
	AbortTrans(tx interface{}) error
	CommitTrans(tx interface{}) error

	//Domain
	GetDomainOfBucketDomain(domainHost string) (info types.DomainInfo, err error)
	GetBucket(bucket string) (uid string, err error)
	GetDomain(projectId string, domainHost string) (info types.DomainInfo, err error)
	GetDomainInfos(projectId string, bucketDomain string) (info []types.DomainInfo, err error)
	GetAllDomainInfos(projectId string) (info []types.DomainInfo, err error)
	InsertDomain(info types.DomainInfo) (err error)
	DelDomain(info types.DomainInfo) (err error)

	//DomainTls
	UpdateDomainTls() error
}
