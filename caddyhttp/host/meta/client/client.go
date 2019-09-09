package client

import "github.com/journeymidnight/yig-front-caddy/caddyhttp/host/meta/types"

type Client interface {
	//Transaction
	NewTrans() (tx interface{}, err error)
	AbortTrans(tx interface{}) error
	CommitTrans(tx interface{}) error

	//Domain
	GetDomainOfBucketDomain(projectId string, domainHost string) (info types.DomainInfo, err error)
	GetBucket(bucket string) (uid string, err error)
	GetDomain(projectId string, domainHost string) (info types.DomainInfo, err error)
	GetDomainInfos(projectId string) (info []types.DomainInfo, err error)
	InsertDomain(info types.DomainInfo, tx interface{}) (err error)
	DelDomain(info types.DomainInfo, tx interface{}) (err error)

	//DomainTls
	UpdateDomainTls(info types.DomainTlsInfo) error
}
