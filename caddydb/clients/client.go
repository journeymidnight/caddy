package clients

import (
	"github.com/journeymidnight/yig-front-caddy/caddydb/types"
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
	GetDomainInfos(projectId string, bucketDomain string, sealKey string) (info []types.DomainInfo, err error)
	GetAllDomainInfos(projectId string, sealKey string) (info []types.DomainInfo, err error)
	InsertDomain(info types.DomainInfo, secretKey string) (err error)
	DelDomain(info types.DomainInfo) (err error)

	//DomainTls
	UpdateDomainTLS(info types.DomainInfo, sealKey string) (err error)

	//ImageStyle
	InsertStyle(style types.ImageStyle) (err error)
	DelStyle(style types.ImageStyle) (err error)
	UpdateStyle(style types.ImageStyle) (err error)
	GetStyles(bucket string) (styles []types.ImageStyle, err error)
	GetStyle(bucket string, styleName string) (style types.ImageStyle, err error)
}
