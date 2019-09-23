package types

type DomainTlsInfo struct {
	ProjectId string
	Tls       string
	TlsKey    string
}

func (c *DomainTlsInfo) UpdateDomainTlsInfo() (string, []interface{}) {
	sql := ""
	args := []interface{}{}
	return sql, args
}
