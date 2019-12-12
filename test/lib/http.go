package lib

import (
	"io"
	"io/ioutil"
	"net/http"
)

func HttpPut(url string, header map[string]string, payload io.Reader) (status int, body string) {
	req, _ := http.NewRequest("PUT", url, payload)
	for k, v := range header {
		req.Header.Add(k, v)
	}
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	status = res.StatusCode
	bodyByte, _ := ioutil.ReadAll(res.Body)
	body = string(bodyByte)
	return
}

func HttpGet(url string, header map[string]string) (status int, body string) {
	req, _ := http.NewRequest("GET", url, nil)
	for k, v := range header {
		req.Header.Add(k, v)
	}
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	status = res.StatusCode
	bodyByte, _ := ioutil.ReadAll(res.Body)
	body = string(bodyByte)
	return
}

func HttpDelete(url string, header map[string]string) (status int, body string) {
	req, _ := http.NewRequest("DELETE", url, nil)
	for k, v := range header {
		req.Header.Add(k, v)
	}
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	status = res.StatusCode
	bodyByte, _ := ioutil.ReadAll(res.Body)
	body = string(bodyByte)
	return
}

func HttpPost(url string, header map[string]string, payload io.Reader) (status int, body string) {
	req, _ := http.NewRequest("POST", url, payload)
	for k, v := range header {
		req.Header.Add(k, v)
	}
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	status = res.StatusCode
	bodyByte, _ := ioutil.ReadAll(res.Body)
	body = string(bodyByte)
	return
}
