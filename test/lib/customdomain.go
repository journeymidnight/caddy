package lib

import "fmt"

func NewCustomDomain(domainHost string) (int, error) {
	url := "http://" + End_Point + "/?x-oss-action=NewCustomDomain"
	token, err := GetJwt(domainHost)
	if err != nil {
		return 0, err
	}
	status, err := SetRequest("PUT", url, token)
	return status, err
}

func GetCustomDomain(domainHost string) (int, error) {
	url := "http://" + End_Point + "/?x-oss-action=GetCustomDomain"
	token, err := GetJwt(domainHost)
	if err != nil {
		return 0, err
	}
	status, err := SetRequest("GET", url, token)
	return status, err
}

func DelCustomDomain(domainHost string) (int, error) {
	url := "http://" + End_Point + "/?x-oss-action=DelCustomDomain"
	token, err := GetJwt(domainHost)
	if err != nil {
		return 0, err
	}
	status, err := SetRequest("DELETE", url, token)
	return status, err
}

func CustomDomainAccess(domainHost string) (int, error) {
	url := "http://" + domainHost + "/" + TEST_KEY
	fmt.Println("Custom domain access path", url)
	status, err := SetRequestWithDomain("GET", url, "")
	return status, err
}
