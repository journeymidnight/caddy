package tidbclient

import (
	"database/sql"
	"github.com/journeymidnight/yig-front-caddy/caddyhttp/host/meta/types"
)

func (DB *TidbClient) GetDomainOfBucketDomain(projectId string, domainHost string) (info types.DomainInfo, err error) {
	var domainBucket string
	sql, args := info.GetDomainOfBucketDomain(projectId, domainHost)
	row, err := DB.ClientBusiness.Query(sql, args...)
	if err != nil {
		return
	}
	defer row.Close()
	row.Next()
	err = row.Scan(&domainBucket)
	info.DomainBucket = domainBucket
	return
}

func (DB *TidbClient) GetBucket(bucket string) (uid string, err error) {
	sql := "select uid from buckets where bucketname=?"
	args := []interface{}{bucket}
	row, err := DB.ClientS3.Query(sql, args...)
	if err != nil {
		return
	}
	defer row.Close()
	row.Next()
	err = row.Scan(&uid)
	if err != nil {
		return
	}
	return
}

func (DB *TidbClient) GetDomain(projectId string, domainHost string) (info types.DomainInfo, err error) {
	var pid, domainH, domainB string
	sql, args := info.GetDomain(projectId, domainHost)
	row, err := DB.ClientBusiness.Query(sql, args...)
	if err != nil {
		return
	}
	defer row.Close()
	row.Next()
	err = row.Scan(&pid, &domainH, &domainB)
	if err != nil {
		return
	}
	info.ProjectId = pid
	info.DomainHost = domainH
	info.DomainBucket = domainB
	return
}

func (DB *TidbClient) GetDomainInfos(projectId string) (info []types.DomainInfo, err error) {
	sql := "select * from custom_domain where project_id=?"
	args := []interface{}{projectId}
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

func (DB *TidbClient) InsertDomain(info types.DomainInfo, tx interface{}) (err error) {
	var sqlTx *sql.Tx
	tx, err = DB.ClientBusiness.Begin()
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

func (DB *TidbClient) DelDomain(info types.DomainInfo, tx interface{}) (err error) {
	var sqlTx *sql.Tx
	tx, err = DB.ClientBusiness.Begin()
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
