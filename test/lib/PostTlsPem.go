package lib

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
)

var (
	uploadFileKey = "upload-key"
)

func PostTlsPem(domainHost string) error {
	token, err := GetJwt(domainHost)
	if err != nil {
		return err
	}
	url := "http://s3.test.com/?x-oss-action=TlsPutCustomDomain"
	params := map[string]string{
		"tls_domain": `-----BEGIN CERTIFICATE-----
MIIEGjCCAoKgAwIBAgIQCm6NcDHIdLV/ybuSPoAVyTANBgkqhkiG9w0BAQsFADBR
MR4wHAYDVQQKExVta2NlcnQgZGV2ZWxvcG1lbnQgQ0ExEzARBgNVBAsMCnJvb3RA
Ym9nb24xGjAYBgNVBAMMEW1rY2VydCByb290QGJvZ29uMB4XDTE5MDYwMTAwMDAw
MFoXDTI5MDkyNjE0NDczN1owPjEnMCUGA1UEChMebWtjZXJ0IGRldmVsb3BtZW50
IGNlcnRpZmljYXRlMRMwEQYDVQQLDApyb290QGJvZ29uMIIBIjANBgkqhkiG9w0B
AQEFAAOCAQ8AMIIBCgKCAQEA6bssFdooNL5WNYjXsm31kn4LedfUeAd3aTmkJSYQ
5QSZjPN399B2LxodLwAlhiuVKhljjAT7YsLeO2KjudJkuSeW9IG6j7OcgPdfcQq8
+cMYjlb+EjomcXwYNl/Go1+vXWcMRZmWhBLP4aqvmPedBpfs2Is+59yfOQFrdasH
0XFsPGhizyz3x5LZNpRea39juyQicQ1ndZ9njneT02CoYZ2qrp5uxJOwNGmx0XAN
KSBRIffTijLbVARVcbXXRFbsnMZg9AY+wdlFgAofyjJo4eW0D1AB+2sOeJbSE8/Y
L4BPwfM3GdbDTkedJ9D5us6WS+4G4NAF43hj1vFgTO0z5QIDAQABo4GAMH4wDgYD
VR0PAQH/BAQDAgWgMBMGA1UdJQQMMAoGCCsGAQUFBwMBMAwGA1UdEwEB/wQCMAAw
HwYDVR0jBBgwFoAUy8IGMkO0o048QeepTiyEqV0irJMwKAYDVR0RBCEwH4Idd3d3
Lndvaml1c2hpeGlhbmdzaGl5aXNoaS54eXowDQYJKoZIhvcNAQELBQADggGBANOW
9bZ+OwKhJgj0jvyzOHzmHTdvMUEQt4ztTE7NuolOPGAD3v8In3PDcOmKqn8tEMlE
MTRUU3ntEA8wJRNye76StSO/Ns+V+Le5VPRuFUkaga0tPo/tenu2Bj0dPypdB21D
kMndDvoYQR9N+o5ZKIZRUI2W9+212nUOPbAnjEjOgdwAspP5ng4C0cl3E9+AMOMj
cCIvx5i3KmJlovlt+C3P/AnGzMrLQ7VJVFgR2omhMp1O7LxUZ0BXSrGh3gjLtizK
ofY88IQSDy8bNAm4ftbTAW9zVbYJ/jn9hou7JCoDQMYpdRqA1rLkFJv2kgHzhmZA
havtqp+oAre+ykwSjIJuFTM5l/MxPBRPs/foOffXJmYqU3Xd+OugZw/g6uKsHD0X
fmZtJHxFJWa2UiUxD9+fbCJ3EIrvtRZOt42kTCdk5RSdhVJM5qcGIfHqe1Qp/eeK
4PrKnpJtfK/r+yLz/aaUdkW6NClx1vX2CnPh4/AIzpYWdOcoIBLHKr/RkeJ4Rg==
-----END CERTIFICATE-----
`,
		"tls_domain_key": `-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEA9MxefCqZrX2/fKbhg097pTBlQh49j0r/UgjzmAvgCm5HKO6N
20k1++zXsJdbvfANJW0xKkQcx5x5dmd6HtUbl5jBhP5WKIq0uKf1CT9GxV+KDsem
FDmVDZt4u+qUh2WyF8re7lCQtPHjTxIMExdUUuRq3+31sxZeXbZTHk4xbIsIw57s
GMgqmIXZikyFS4BEdSqKrODdlHs4yHGb6MRNtbAvkyoO6IMQaMESx9pL58s3V078
Vli27NwtKw7aXwmWRX5aKuOlbOyHGvkM1x6g+h2ASQdcm37x/adoc1FJX3a+0wnj
zY6ythQbOzTQr70GRp7UURRpwlXb1nBKX4eM2wIDAQABAoIBAQDDrfW9JURw1TrT
BobvswJwWk4m2wt0VovfZH5rxNpTHDHPu1kzt8LqXTlQ1LymCJRJnY4WbhnfSeoh
xrdsnAcRWC8hAzmV8MJjdQa/XJBSx3sRM3BDbIAEP4n0FKCj1pvUvvsS3t29djmw
hAmY3SYvmJ0B4TJs9G18VXj9mZ2SV8Tj3YvntJdEtABpUhfx9urWqLQ0gaZqrfFA
MkTZupmfq84H0iZCPXJ3v85mFppg8TmLz+GmtD97D8mM5jIS5KHW5rHY/MrB+IMo
pXxnHS36kz+Lfb/IqJfeq4bVh2lKrgWtGZp6495kg9FmeFZsSgVJr9DTK/sXBaZp
NQobIZF5AoGBAP6euH2uvYnCCHrwOqOL5lMybfEYPUQRqQs2T+hfNl3r71LoXGqb
rNPx/iQ8lOvHoKo6rJ6tkXi3m8jzwzagLY/wvktpvTFVzFGnlEfLjF+fwdctCNe2
TlkI5w5v5flVjYdSMsaTI9Q+kAy7ZtHr18JBurGm5PSJMH/+B5KwzEt9AoGBAPYg
BWPFDdWGYC5ATG2dk2ts80mpkI04Jn0dsYlV8Zbvx5ZprqAevzJgdQl15YK5y+cu
+NL5/bHXNx0OqeoFV4cfXr609apzhHqBCf+1C130wKORsGxcrENJL9bjVjU0kggD
fq50MPmMEDy3YKpVoTSoHxJN3Jt8c6vJp2SVTrk3AoGBAMR5n61ECHfjpP2Qn/8R
diKe2F1hMe1znzjKqTCLP/2LkLouGRoelGdb5Zr6sIOcIGsbKZasIoO/UtPm3a0q
Pt5Ci+TPbDFEnownvvEDrYcgVMiAtMEXqS9lAj6OtCwJC1PvZsT6R3yA5lEczsOP
tIbuh2yMd0IHFsi9MgTMrmhVAoGALBIpxOapBn3sRPvgc5ROrGKy4ZLKPm86fRbP
Dm9kNgmzbFx6F1PKGqQo8Mu6kADi4P+JMIXxBmIqDTrt8+iG9rwIRA+1GZNbum/W
sYHiii1kPSW+OHkTo2y8czb73cUPDP5LNcO6bUTGN4kCdx2kIwCjp6xnfzP2pmES
ZzFRClsCgYEAq11XI1PHvanA0+6C1pPI4tqQB2PPH0iq8fBedfJJqTvQf7+B6MBU
XeyIg8Le2uOstjwpIV9vIXGd6h5UAcVRqQbKPfz0q4RoH5j3wGO7M3M/+HDRMrcm
nh+N76wRR+58+yUWhztWQoVPnB4j0AgGTdtTfapmOgF/aicvXn3AyfA=
-----END RSA PRIVATE KEY-----

`,
	}
	req, err := NewFileUploadRequest(url, params)
	if err != nil {
		return fmt.Errorf("error to new upload file request:%s\n", err.Error())
	}
	req.Header.Set("Authorization", token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error to request to the server:%s\n", err.Error())
	}
	body := &bytes.Buffer{}
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

// NewFileUploadRequest ...
func NewFileUploadRequest(url string, params map[string]string) (*http.Request, error) {
	body := &bytes.Buffer{}
	// 文件写入 body
	writer := multipart.NewWriter(body)
	// 其他参数列表写入 body
	for k, v := range params {
		if err := writer.WriteField(k, v); err != nil {
			return nil, err
		}
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", writer.FormDataContentType())
	return req, err
}

func DelTlsPem(domainHost string) (int, error) {
	url := "http://" + End_Point + "/?x-oss-action=DelCustomDomain"
	token, err := GetJwt(domainHost)
	if err != nil {
		return 0, err
	}
	status, err := SetRequest("DELETE", url, token)
	return status, err
}
