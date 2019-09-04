package host

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type TidbClient struct {
	ClientS3       *sql.DB
	ClientBusiness *sql.DB
}

func GetDatabaseConn(h Host) (*TidbClient, error) {
	cli := &TidbClient{}
	s3Conn, err := sql.Open("mysql", h.S3Source)
	cli.ClientS3 = s3Conn
	bsConn, err := sql.Open("mysql", h.BusinessSource)
	cli.ClientBusiness = bsConn
	return cli, err
}

func GetTargetBucket(conn *TidbClient, host string, projectID string) (string, error) {
	var bucket_domain string
	sql := "select bucket_domain from custom_domain where host=? and project_id=?"
	args := []interface{}{host, projectID}
	rows, err := conn.ClientBusiness.Query(sql, args...)
	if err != nil {
		return bucket_domain, fmt.Errorf("Get the target bucket failure:", err)
	}
	defer rows.Close()
	rows.Next()
	err = rows.Scan(&bucket_domain)
	if rows == nil {
		return bucket_domain, fmt.Errorf("No host to router:", bucket_domain)
	}
	return bucket_domain, err
}

func AuthBucketInfo(domain Domain) ([]string, error) {
	var args []interface{}
	var bucketname string
	var result []string
	custom := domain.customInfo
	sql := "select bucketname from users where userid=? and bucketname=?"
	if domain.editBucket != "" {
		args = []interface{}{custom.projectId, domain.editBucket}
	} else {
		args = []interface{}{custom.projectId, custom.bucket}
	}
	rows, err := domain.conn.ClientS3.Query(sql, args...)
	for rows.Next() {
		err = rows.Scan(&bucketname)
		if err != nil {
			return result, err
		}
		result = append(result, bucketname)
	}
	defer rows.Close()
	return result, nil
}

func AuthHostInfo(domain Domain) ([]string, error) {
	var host string
	var result []string
	custom := domain.customInfo
	sql := "select host from custom_domain where host=? and project_id=?"
	args := []interface{}{custom.domainHost, custom.projectId}
	rows, err := domain.conn.ClientBusiness.Query(sql, args...)
	for rows.Next() {
		err = rows.Scan(&host)
		if err != nil {
			return result, err
		}
		result = append(result, host)
	}
	defer rows.Close()
	return result, nil
}

func GetDomainInfo(domain Domain) (customDomainInfo CustomDomainInfo, err error) {
	custom := domain.customInfo
	sql := "select * from custom_domain where project_id=?"
	args := []interface{}{custom.projectId}
	rows, err := domain.conn.ClientBusiness.Query(sql, args...)
	if err != nil {
		return
	}
	for rows.Next() {
		IcustomDomain := CustomDomain{}
		err = rows.Scan(&IcustomDomain.ProjectId, &IcustomDomain.DomainHost, &IcustomDomain.Domain, &IcustomDomain.TlsKey, &IcustomDomain.Tls)
		customDomainInfo.CustomDomain = append(customDomainInfo.CustomDomain, IcustomDomain)
	}
	if err != nil {
		return
	}
	defer rows.Close()
	return
}

func InsertDomainInfo(domain Domain) error {
	custom := domain.customInfo
	var sqlTx *sql.Tx
	var tx interface{}
	tx, err := domain.conn.ClientBusiness.Begin()
	defer func() {
		if err == nil {
			err = sqlTx.Commit()
		}
		if err != nil {
			sqlTx.Rollback()
		}
	}()
	sqlTx, _ = tx.(*sql.Tx)
	sql := "insert into custom_domain values(?,?,?,?,?)"
	args := []interface{}{custom.projectId, custom.domainHost, custom.domain, "", ""}
	_, err = sqlTx.Exec(sql, args...)
	if err != nil {
		return err
	}
	return nil
}

func EditDomainInfo(domain Domain, editBucket string) error {
	custom := domain.customInfo
	sql := "update custom_domain set bucket_domain=? where host=? and project_id=?"
	args := []interface{}{editBucket, custom.domainHost, custom.projectId}
	_, err := domain.conn.ClientBusiness.Exec(sql, args...)
	if err != nil {
		return err
	}
	return nil
}

func DelDomainInfo(domain Domain) error {
	custom := domain.customInfo
	var sqlTx *sql.Tx
	var tx interface{}
	tx, err := domain.conn.ClientBusiness.Begin()
	defer func() {
		if err == nil {
			err = sqlTx.Commit()
		}
		if err != nil {
			sqlTx.Rollback()
		}
	}()
	sqlTx, _ = tx.(*sql.Tx)
	sql := "delete from custom_domain where host=? and project_id=?"
	args := []interface{}{custom.domainHost, custom.projectId}
	_, err = sqlTx.Exec(sql, args...)
	if err != nil {
		return err
	}
	return nil
}

func UpdateDomainInfoTls(domain Domain) error {
	tls := domain.tlsInfo
	custom := domain.customInfo
	sql := "update custom_domain set tls_key=? and tls=? where host=? and project_id=?"
	args := []interface{}{tls.tlsKey, tls.tls, custom.domainHost, custom.projectId}
	_, err := domain.conn.ClientBusiness.Exec(sql, args)
	if err != nil {
		return err
	}
	return nil
}
