package host

import (
	"bytes"
	"encoding/xml"
	"time"
)

const (
	timeFormatAMZ = "2006-01-02T15:04:05.000Z" // Reply date format
)

type CustomDomainResponse struct {
	XMLName      xml.Name `xml:"http://www.unicloud.com CustomDomainResult" json:"-"`
	LastModified string   // time string of format "2006-01-02T15:04:05.000Z"
}

func GenerateObjectResponse(lastModified time.Time) CustomDomainResponse {
	lastModified = time.Now()
	return CustomDomainResponse{
		LastModified: lastModified.UTC().Format(timeFormatAMZ),
	}
}

func EncodeResponse(response interface{}) []byte {
	var bytesBuffer bytes.Buffer
	bytesBuffer.WriteString(xml.Header)
	e := xml.NewEncoder(&bytesBuffer)
	e.Encode(response)
	return bytesBuffer.Bytes()
}
