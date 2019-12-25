package pipa

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/journeymidnight/yig-front-caddy/caddyerrors"
	"strconv"
	"strings"
)

type ErrorResponse struct {
	XMLName        xml.Name `xml:"Error" json:"-"`
	CaddyErrorCode string   `xml:"error"`
	Description    string   `xml:"description"`
	HttpStatusCode int      `xml:"status"`
}

func writeErrorResponse(err caddyerrors.HandleError) (result []byte, error error) {
	errResponse := new(ErrorResponse)
	errResponse.CaddyErrorCode = err.CaddyErrorCode()
	errResponse.Description = err.Description()
	errResponse.HttpStatusCode = err.HttpStatusCode()
	var bytesBuffer bytes.Buffer
	bytesBuffer.WriteString(xml.Header)
	e := xml.NewEncoder(&bytesBuffer)
	e.Encode(errResponse)
	error = err
	return bytesBuffer.Bytes(), error
}

type ErrorResponseFromPipa struct {
	XMLName        xml.Name `xml:"Error" json:"-"`
	Description    string   `xml:"description"`
	HttpStatusCode int      `xml:"status"`
}

func writeErrorResponseWithPipa(err string) (result []byte, error error) {
	errAll := strings.Split(err, ",")
	errResponse := new(ErrorResponseFromPipa)
	errResponse.Description = errAll[1]
	code, _ := strconv.Atoi(errAll[0])
	errResponse.HttpStatusCode = code
	var bytesBuffer bytes.Buffer
	bytesBuffer.WriteString(xml.Header)
	e := xml.NewEncoder(&bytesBuffer)
	e.Encode(errResponse)
	error = fmt.Errorf("Pipa image process err! ")
	return bytesBuffer.Bytes(), error
}
