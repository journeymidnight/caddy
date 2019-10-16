package lib

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

func PostTlsPem(domainHost string) error {
	token, err := GetJwt(domainHost)
	if err != nil {
		return err
	}

	pemName := FileDir + domainHost + ".pem"
	pemKey := FileDir + domainHost + "-key.pem"
	filePem, err := ReadPemAndKey(pemName)
	if err != nil {
		return err
	}
	filePemKey, err := ReadPemAndKey(pemKey)
	if err != nil {
		return err
	}

	url := "http://s3.test.com/?x-oss-action=PutCertificate"
	params := map[string]string{
		"tls_domain":     string(filePem),
		"tls_domain_key": string(filePemKey),
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
	url := "https://" + End_Point + "/?x-oss-action=DelCertificate"
	token, err := GetJwt(domainHost)
	if err != nil {
		return 0, err
	}
	status, err := SetRequest("DELETE", url, token)
	return status, err
}

func ReadPemAndKey(filePth string) ([]byte, error) {
	f, err := os.Open(filePth)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(f)
}
