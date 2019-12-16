package tidbclient

import (
	. "github.com/journeymidnight/yig-front-caddy/caddyerrors"
)

func (DB *TidbClient) GetBucket(bucket string) (uid string, err error) {
	sql := "select uid from buckets where bucketname=?"
	args := []interface{}{bucket}
	row := DB.Client.QueryRow(sql, args...)
	err = row.Scan(&uid)
	if err != nil {
		return uid, ErrNoSuchKey
	}
	return uid, nil
}
