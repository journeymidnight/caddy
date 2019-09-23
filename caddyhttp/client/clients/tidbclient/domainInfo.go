package tidbclient

import (
	"database/sql"
	"fmt"
	"github.com/journeymidnight/yig-front-caddy/caddyhttp/client/types"
)

func (DB *TidbClient) GetDomainOfBucketDomain(domainHost string) (info types.DomainInfo, err error) {
	var domainBucket string
	sql := info.GetDomainOfBucketDomain()
	args := []interface{}{domainHost}
	row := DB.ClientBusiness.QueryRow(sql, args...)
	err = row.Scan(&domainBucket)
	info.DomainBucket = domainBucket
	return
}

func (DB *TidbClient) GetBucket(bucket string) (uid string, err error) {
	sql := "select uid from buckets where bucketname=?"
	args := []interface{}{bucket}
	row := DB.ClientS3.QueryRow(sql, args...)
	err = row.Scan(&uid)
	if err != nil {
		return
	}
	return uid, nil
}

func (DB *TidbClient) GetDomain(projectId string, domainHost string) (info types.DomainInfo, err error) {
	var pid, domainH, domainB string
	sql := info.GetDomain()
	args := []interface{}{projectId, domainHost}
	row := DB.ClientBusiness.QueryRow(sql, args...)
	err = row.Scan(&pid, &domainH, &domainB)
	if err != nil {
		return info, fmt.Errorf("No such key!")
	}
	info.ProjectId = pid
	info.DomainHost = domainH
	info.DomainBucket = domainB
	return
}

func (DB *TidbClient) GetDomainInfos(projectId string, bucketDomain string) (info []types.DomainInfo, err error) {
	sql := "select * from custom_domain where project_id=? and bucket_domain=?"
	args := []interface{}{projectId, bucketDomain}
	rows, err := DB.ClientBusiness.Query(sql, args...)
	if err != nil {
		return
	}
	for rows.Next() {
		ICustomDomain := types.DomainInfo{}
		err = rows.Scan(&ICustomDomain.ProjectId, &ICustomDomain.DomainHost, &ICustomDomain.DomainBucket)
		info = append(info, ICustomDomain)
	}
	if err != nil {
		return
	}
	defer rows.Close()
	return
}

func (DB *TidbClient) InsertDomain(info types.DomainInfo) (err error) {
	var sqlTx *sql.Tx
	var tx interface{}
	tx, err = DB.ClientBusiness.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err == nil {
			err = sqlTx.Commit()
		}
		if err != nil {
			sqlTx.Rollback()
		}
	}()
	sqlTx, _ = tx.(*sql.Tx)
	sql, args := info.InsertDomain()
	_, err = sqlTx.Exec(sql, args...)
	if err != nil {
		return err
	}
	return nil
}

func (DB *TidbClient) DelDomain(info types.DomainInfo) (err error) {
	var sqlTx *sql.Tx
	var tx interface{}
	tx, err = DB.ClientBusiness.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err == nil {
			err = sqlTx.Commit()
		}
		if err != nil {
			sqlTx.Rollback()
		}
	}()
	sqlTx, _ = tx.(*sql.Tx)
	sql, args := info.DeleteDomain()
	_, err = sqlTx.Exec(sql, args...)
	if err != nil {
		return err
	}
	return nil
}
