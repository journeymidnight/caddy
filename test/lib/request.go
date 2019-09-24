package lib

import (
	"net/http"
)

func SetRequest(method string, url string, token string) (int, error) {
	client := &http.Client{}
	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	request.Header.Set("Authorization", token)
	response, _ := client.Do(request)
	status := response.StatusCode
	return status, err
}

func SetRequestWithDomain(method string, url string, token string) (int, error) {
	client := &http.Client{}
	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	response, _ := client.Do(request)
	status := response.StatusCode
	return status, err
}
