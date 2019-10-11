package types

type DomainInfo struct {
	ProjectId    string
	Bucket       string
	DomainBucket string
	DomainHost   string
	TlsDomain    string
	TlsDomainKey string
}

func (c *DomainInfo) GetDomainOfBucketDomain() string {
	sql := "select bucket_domain from custom_domain where host_domain=?"
	return sql
}

func (c *DomainInfo) GetDomain() string {
	sql := "select project_id,host_domain,bucket_domain from custom_domain where project_id=? and host_domain=?"
	return sql
}

func (c *DomainInfo) InsertDomain(secretKey string) (string, []interface{}) {
	sql := "insert into custom_domain(project_id,host_domain,bucket_domain,tls_domain,tls_domain_key) values(?,?,?,AES_ENCRYPT(?, ?),AES_ENCRYPT(?, ?))"
	args := []interface{}{c.ProjectId, c.DomainHost, c.DomainBucket, c.TlsDomain, secretKey, c.TlsDomainKey, secretKey}
	return sql, args
}

func (c *DomainInfo) DeleteDomain() (string, []interface{}) {
	sql := "delete from custom_domain where project_id=? and host_domain=?"
	args := []interface{}{c.ProjectId, c.DomainHost}
	return sql, args
}

func (c *DomainInfo) UpdateDomainTls(secretKey string) (string, []interface{}) {
	sql := "update custom_domain set tls_domain=AES_ENCRYPT(?, ?), tls_domain_key= AES_ENCRYPT(?, ?)  where project_id=? and host_domain=?"
	args := []interface{}{c.TlsDomain, secretKey, c.TlsDomainKey, secretKey, c.ProjectId, c.DomainHost}
	return sql, args
}
