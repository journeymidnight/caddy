package types

type DomainInfo struct {
	ProjectId    string
	Bucket       string
	DomainBucket string
	DomainHost   string
}

func (c *DomainInfo) GetDomainOfBucketDomain() string {
	sql := "select bucket_domain from custom_domain where host_domain=?"
	return sql
}

func (c *DomainInfo) GetDomain() string {
	sql := "select * from custom_domain where project_id=? and host_domain=?"
	return sql
}

func (c *DomainInfo) InsertDomain() (string, []interface{}) {
	sql := "insert into custom_domain values(?,?,?)"
	args := []interface{}{c.ProjectId, c.DomainHost, c.DomainBucket}
	return sql, args
}

func (c *DomainInfo) DeleteDomain() (string, []interface{}) {
	sql := "delete from custom_domain where project_id=? and host_domain=?"
	args := []interface{}{c.ProjectId, c.DomainHost}
	return sql, args
}
