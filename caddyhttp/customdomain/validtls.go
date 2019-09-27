package customdomain

import (
	"crypto/tls"
	"crypto/x509"
	. "encoding/pem"
	. "github.com/journeymidnight/yig-front-caddy/caddyerrors"
	"net/http"
	"strings"
)

func GetTlsFromFormData(r *http.Request) (tls string, tlsKey string, err error) {
	DOMAIN.Log.Println(20, "Enter get tls from form data")
	r.ParseMultipartForm(128 << 10)
	if r.MultipartForm != nil {
		tls = r.MultipartForm.Value["tls_domain"][0]
		tlsKey = r.MultipartForm.Value["tls_domain_key"][0]
		if tls == "" || tlsKey == "" {
			return "", "", ErrFormDataParameterParsing
		}
		DOMAIN.Log.Println(20, "Get tls from form data succeed!")
		return tls, tlsKey, nil
	}
	return "", "", ErrFormDataParameterParsing
}

func ValidTls(hostDomain, tlsCrt, tlsKey string) error {
	DOMAIN.Log.Println(20, "Enter valid tls and tls key.")
	resultCrt, _ := Decode([]byte(tlsCrt))
	if resultCrt.Type != "CERTIFICATE" {
		return ErrInvalidTlsPem
	}
	pub, err := x509.ParseCertificate(resultCrt.Bytes)
	if err != nil {
		return ErrInvalidTlsPem
	}
	dnsName := pub.DNSNames
	validKey := false
	for _, host := range dnsName {
		if hostDomain == host {
			validKey = true
		}
		if strings.HasSuffix(hostDomain, host) {
			validKey = true
			break
		}
		ht := strings.Split(host,"*")
		for _, h := range ht{
			if strings.HasSuffix(hostDomain, h) {
				validKey = true
				break
			}			
		}
	}
	DOMAIN.Log.Printf(20, "Got a %T with dns name : %v %v", pub, dnsName, hostDomain)
	if !validKey {
		return ErrInvalidTlsPem
	}
	_, err = tls.X509KeyPair([]byte(tlsCrt), []byte(tlsKey))
	if err != nil {
		return ErrInvalidTls
	}
	return nil
}
