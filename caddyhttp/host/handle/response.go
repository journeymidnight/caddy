package handle

import (
	"bytes"
	"encoding/xml"
	"time"
)

const (
	timeFormatAMZ = "2006-01-02T15:04:05.000Z" // Reply date format
)

type DomainResponse struct {
	ProjectId  string `xml:"project-id"`
	DomainHost string `xml:"host"`
	Domain     string `xml:"bucket-domain"`
	TlsKey     string `xml:"tls-key"`
	Tls        string `xml:"tls"`
}

type Response struct {
	XMLName      xml.Name `xml:"http://www.unicloud.com CustomDomainResult" json:"-"`
	LastModified string   // time string of format "2006-01-02T15:04:05.000Z"
}

func GenerateObjectResponse(lastModified time.Time) Response {
	lastModified = time.Now()
	return Response{
		LastModified: lastModified.UTC().Format(timeFormatAMZ),
	}
}

func EncodeResponse(result []byte, response interface{}) []byte {
	var bytesBuffer bytes.Buffer
	var bingo []byte
	bytesBuffer.WriteString(xml.Header)
	e := xml.NewEncoder(&bytesBuffer)
	e.Encode(response)
	bingo = bytesBuffer.Bytes()
	bingo = append(bingo, result...)
	return bingo
}
