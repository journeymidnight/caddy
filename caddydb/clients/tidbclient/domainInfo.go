package tidbclient

import (
	"database/sql"
	"github.com/journeymidnight/yig-front-caddy/caddydb/types"
	. "github.com/journeymidnight/yig-front-caddy/caddyerrors"
)

func (DB *TidbClient) GetDomainOfBucketDomain(domainHost string) (info types.DomainInfo, err error) {
	var domainBucket string
	sql := info.GetDomainOfBucketDomain()
	args := []interface{}{domainHost}
	row := DB.Client.QueryRow(sql, args...)
	err = row.Scan(&domainBucket)
	info.DomainBucket = domainBucket
	return
}

func (DB *TidbClient) GetDomain(projectId string, domainHost string) (info types.DomainInfo, err error) {
	var pid, domainH, domainB string
	sql := info.GetDomain()
	args := []interface{}{projectId, domainHost}
	row := DB.Client.QueryRow(sql, args...)
	err = row.Scan(&pid, &domainH, &domainB)
	if err != nil {
		return info, ErrNoSuchKey
	}
	info.ProjectId = pid
	info.DomainHost = domainH
	info.DomainBucket = domainB
	return
}

func (DB *TidbClient) GetDomainInfos(projectId string, bucketDomain string, sealKey string) (info []types.DomainInfo, err error) {
	sql := "select project_id,host_domain,bucket_domain,IFNULL(AES_DECRYPT(tls_domain, ?),'') from custom_domain where project_id=? and bucket_domain=?"
	args := []interface{}{sealKey, projectId, bucketDomain}
	rows, err := DB.Client.Query(sql, args...)
	if err != nil {
		return info, ErrInvalidSql
	}
	for rows.Next() {
		ICustomDomain := types.DomainInfo{}
		err = rows.Scan(&ICustomDomain.ProjectId, &ICustomDomain.DomainHost, &ICustomDomain.DomainBucket, &ICustomDomain.TlsDomain)
		info = append(info, ICustomDomain)
	}
	if err != nil {
		return info, ErrNoSuchKey
	}
	defer rows.Close()
	return
}

func (DB *TidbClient) GetAllDomainInfos(projectId string, sealKey string) (info []types.DomainInfo, err error) {
	sql := "select project_id,host_domain,bucket_domain,IFNULL(AES_DECRYPT(tls_domain, ?),'') from custom_domain where project_id=?"
	args := []interface{}{sealKey, projectId}
	rows, err := DB.Client.Query(sql, args...)
	if err != nil {
		return info, ErrInvalidSql
	}
	for rows.Next() {
		ICustomDomain := types.DomainInfo{}
		err = rows.Scan(&ICustomDomain.ProjectId, &ICustomDomain.DomainHost, &ICustomDomain.DomainBucket, &ICustomDomain.TlsDomain)
		info = append(info, ICustomDomain)
	}
	if err != nil {
		return info, ErrNoSuchKey
	}
	defer rows.Close()
	return
}

func (DB *TidbClient) InsertDomain(info types.DomainInfo, sealKey string) (err error) {
	var sqlTx *sql.Tx
	var tx interface{}
	tx, err = DB.Client.Begin()
	if err != nil {
		return ErrSqlTransaction
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
	sql, args := info.InsertDomain(sealKey)
	_, err = sqlTx.Exec(sql, args...)
	if err != nil {
		return ErrSqlInsert
	}
	return nil
}

func (DB *TidbClient) DelDomain(info types.DomainInfo) (err error) {
	var sqlTx *sql.Tx
	var tx interface{}
	tx, err = DB.Client.Begin()
	if err != nil {
		return ErrSqlTransaction
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
		return ErrSqlDelete
	}
	return nil
}

func (DB *TidbClient) UpdateDomainTLS(info types.DomainInfo, sealKey string) (err error) {
	var sqlTx *sql.Tx
	var tx interface{}
	tx, err = DB.Client.Begin()
	if err != nil {
		return ErrSqlTransaction
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
	sql, args := info.UpdateDomainTls(sealKey)
	_, err = sqlTx.Exec(sql, args...)
	if err != nil {
		return ErrSqlUpdate
	}
	return nil
}

func (DB *TidbClient) GetAllCertificateInfos(sealKey string) (info []types.DomainInfo, err error) {
	sql := "select IFNULL(AES_DECRYPT(tls_domain, ?),''),IFNULL(AES_DECRYPT(tls_domain_key, ?),'') from custom_domain"
	args := []interface{}{sealKey, sealKey}
	rows, err := DB.Client.Query(sql, args...)
	if err != nil {
		return info, ErrInvalidSql
	}
	for rows.Next() {
		ICustomDomain := types.DomainInfo{}
		err = rows.Scan(&ICustomDomain.TlsDomain, &ICustomDomain.TlsDomainKey)
		if err != nil {
			return info, ErrNoSuchKey
		}
		info = append(info, ICustomDomain)
	}
	defer rows.Close()
	return
}
