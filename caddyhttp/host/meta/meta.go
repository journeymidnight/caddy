package meta

import (
	"github.com/journeymidnight/yig-front-caddy/caddyhttp/host/meta/client"
	"github.com/journeymidnight/yig-front-caddy/caddyhttp/host/meta/client/tidbclient"
)

type Meta struct {
	Client client.Client
}

func New(s3Source string, businessSource string) *Meta {
	meta := Meta{}
	meta.Client = tidbclient.NewTidbClient(s3Source, businessSource)
	return &meta
}
