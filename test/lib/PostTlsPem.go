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
MIIEGDCCAoCgAwIBAgIRAKNl6oB3KJ97B7a8ANGw9J0wDQYJKoZIhvcNAQELBQAw
UTEeMBwGA1UEChMVbWtjZXJ0IGRldmVsb3BtZW50IENBMRMwEQYDVQQLDApyb290
QGJvZ29uMRowGAYDVQQDDBFta2NlcnQgcm9vdEBib2dvbjAeFw0xOTA2MDEwMDAw
MDBaFw0yOTA5MjcxMjUzMTZaMD4xJzAlBgNVBAoTHm1rY2VydCBkZXZlbG9wbWVu
dCBjZXJ0aWZpY2F0ZTETMBEGA1UECwwKcm9vdEBib2dvbjCCASIwDQYJKoZIhvcN
AQEBBQADggEPADCCAQoCggEBAKNCLc1/L4Bb8Ri6YGtNk01pAt30e2S5dpKFVmN2
X93ccFfd+kUMgxn8jP8xxIFgifbCf8j6+2nJ2+/IP0xEvHUxLRbZpGT8dTFL8Bz+
uKev7FnS2XRnR1uTfAl4qArU7FORYm587mDiNwJ1meOZY58UnkiNOODpkHlJtS+v
MugQYHmMHp5qq2QHNK0bKR1eKbEnjcL6qL8l0+K6/2IzfsF9OohesgQwWfdx3Sot
6PEMHu/VATVxAjxLWUk+utQXIpo/699RXvJp+R8XTKkw6RkqJBO4YrmzfiKtuFRn
5qKu6ogBEKAuoKLA5K3harjGbr8oSluN95fcLGJKJ9+MiOMCAwEAAaN+MHwwDgYD
VR0PAQH/BAQDAgWgMBMGA1UdJQQMMAoGCCsGAQUFBwMBMAwGA1UdEwEB/wQCMAAw
HwYDVR0jBBgwFoAUy8IGMkO0o048QeepTiyEqV0irJMwJgYDVR0RBB8wHYIbKi53
b2ppdXNoaXhpYW5nc2hpeWlzaGkueHl6MA0GCSqGSIb3DQEBCwUAA4IBgQCaRTPe
ICyJyYg9aAW8LGx0u2YkjbGu/gOZO2TEJfaoT4eTcKL7q2hJAQ5z9tL7sC5hEqUk
cOMwfA552nJPM2ggIaVA7wHbiPp/hDOU5ggxV4jhSrfjbfrYkE5U3wjETgvPCrEq
KOYGNKX25NblwNwops6waknUnWNd9imH8jjZIFsffC/wEbBaeEqh+7iQe2kSEwsZ
A9R04ROiuSuFmDVLvyJGIgbkm1Hm78coo2ioc181ejhikmVxucXyp+IZ6aOLYO75
lX6IcAavy73NblygtXGkxRhHNIiTSjUXLMWbtw1QgzPNaHET20SAkeS2z/homsnx
/4J9NuJq9h5GNLe2lJvlyvXZeigZPS1grBBIuEwM7UrTgU+TfUFEVgI/j2vrBVgr
kbr0wxccadr/nLA4s1RADcvzhVONDQ1S45iTBHQTcw5ecNLefheY6SxbMQlVDMu/
TMykv5CSiUSvFXTPcv205PAOWXU540ZmQ+IF1OUo5rR+O/8juNKMydw95Rs=
-----END CERTIFICATE-----`,
		"tls_domain_key": `-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQCjQi3Nfy+AW/EY
umBrTZNNaQLd9HtkuXaShVZjdl/d3HBX3fpFDIMZ/Iz/McSBYIn2wn/I+vtpydvv
yD9MRLx1MS0W2aRk/HUxS/Ac/rinr+xZ0tl0Z0dbk3wJeKgK1OxTkWJufO5g4jcC
dZnjmWOfFJ5IjTjg6ZB5SbUvrzLoEGB5jB6eaqtkBzStGykdXimxJ43C+qi/JdPi
uv9iM37BfTqIXrIEMFn3cd0qLejxDB7v1QE1cQI8S1lJPrrUFyKaP+vfUV7yafkf
F0ypMOkZKiQTuGK5s34irbhUZ+airuqIARCgLqCiwOSt4Wq4xm6/KEpbjfeX3Cxi
SiffjIjjAgMBAAECggEBAJW8AN5VJHAJ45rWddB0mgGJqsN/lBzRrVq9SHp4+3w0
ziUIxp3F2AhhhcAZGS3AKUlVInZnG2fvZ/UIKGX/NQoXTE+E1i55DdNk3pj5VssV
oyTPKfqFLwFix4S4UQd+OjJ7SVgk5U2erO1ohvgkmyYwduv7+oruLT/QxokxmvMT
ewltA1FYDGv5resAGx8tp2eelctKjck8d9E9UpgR1fReTpEGDFXESJ/0Y5du5UoX
cbA6oBIWHvrWQ2yqvjwVpa+/sF8P/YXOpUOKqqjfG4ONhAk53md0RsgigbGW4EeS
NjevuFWU8jyRBumt22lHUH+RssBVSI02iaZTVOaVGGECgYEA0YgE/r+j34J8bHdG
vGUar4XAYf8c6EzOyEOZPh/Dc03xMlUHS3E7Vqr5f619YfA/3s7jjJ8TWzWC3u/A
zAw098yJl3hJmTdunFjv62fxNDlxNPvtE1pvUdbidahdP/vjGI5GJYlie6YYA9dK
n1ZZdMszJbDaxwZPz5N6VsM0NX0CgYEAx3cOSPN2ARE75ba3Lo73Cu+h2p1YqXLt
OvUXNdkV0D3qozJrOYx2tv5SXBuf9xzjGMKHQG/n/02P+hcPL7XB9ixkeKBWgqKN
MczH+ArhI1dPjQ5dewNCrm6toIXg56IrNtZsXX9Oz7cfMg3rL6l8/+uZwJJWnMDp
iYps0uzohd8CgYBz5LEmmqcwZEMf12VnpOB6vxcm3O2HS5yARmuHYhhAOZc0SLWN
M+cnS9BOn44fUrxxJ6vSxtX0+AcX+jKAaiwN97MO9bh3p6Jllge2BDr0sOT98m4x
6y5xbNK7U1Gop1D37xG7h17Sl47m6PjcYu193TrAGS8ZMFOKs77SKIxDHQKBgQCg
RS5qOY2A0AszybuoomEoHWIM2c8q4Fhzvgk3UEXxvD5zgQidBhtBcFpW/i9rjH+B
HpU0lnZwMi9UwQCH0mCWYBcewZ6heuE+uY+X444BKp+V9IvyUq1aoT3LtKcBF9Hn
TyVlfuyhhD+BpaNq+aGhtPomvK7xZyR/SoWkeY1gOwKBgHboq3zQyr/5ecsCmXZP
/nhWsHt/ZVsKuN6NjPv+mTD4N9lqqpVtVPX5R4UkWgOKm/UfhqQYX72Xba+wKVQ3
9KWc3aEQdqaU+03vL85RJ9IW4tFIudS2kSW/0ZeGwWzvum4BQskvPqnGM14cuqLW
aJRHMTdZ/zG1O/kcvw86A3+a
-----END PRIVATE KEY-----`,
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
	if resp.StatusCode > 300 {
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
	writer := multipart.NewWriter(body)
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
