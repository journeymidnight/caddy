package tidbclient

import (
	"github.com/journeymidnight/yig-front-caddy/caddyhttp/host/meta/types"
)

func (DB *TidbClient) UpdateDomainTls(info types.DomainTlsInfo) error {
	sql, args := info.UpdateDomainTlsInfo()
	_, err := DB.ClientBusiness.Exec(sql, args)
	if err != nil {
		return err
	}
	return nil
}
