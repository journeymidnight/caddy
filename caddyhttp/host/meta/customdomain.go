package meta

import (
	"github.com/journeymidnight/yig-front-caddy/caddyhttp/host/meta/types"
)

func (m *Meta) GetDomainOfBucketDomain(domainHost string) (info types.DomainInfo, err error) {
	info, err = m.Client.GetDomainOfBucketDomain(domainHost)
	return
}

func (m *Meta) ValidBucket(bucket string) (uid string, err error) {
	uid, err = m.Client.GetBucket(bucket)
	return
}

func (m *Meta) GetDomain(projectId string, domainHost string) (info types.DomainInfo, err error) {
	info, err = m.Client.GetDomain(projectId, domainHost)
	return
}

func (m *Meta) GetDomainInfos(projectId string) (info []types.DomainInfo, err error) {
	info, err = m.Client.GetDomainInfos(projectId)
	return
}

func (m *Meta) InsertDomain(customDomainInfo types.DomainInfo) error {
	var tx interface{}
	tx, err := m.Client.NewTrans()
	defer func() {
		if err != nil {
			m.Client.AbortTrans(tx)
		}
	}()
	err = m.Client.InsertDomain(customDomainInfo, tx)
	return err
}

func (m *Meta) DelDomain(customDomainInfo types.DomainInfo) error {
	var tx interface{}
	tx, err := m.Client.NewTrans()
	defer func() {
		if err != nil {
			m.Client.AbortTrans(tx)
		}
	}()
	err = m.Client.DelDomain(customDomainInfo, tx)
	return err
}