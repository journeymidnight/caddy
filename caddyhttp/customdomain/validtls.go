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
	DOMAIN.Log.Info("Enter get tls from form data")
	r.ParseMultipartForm(128 << 10)
	if r.MultipartForm != nil {
		tls = r.MultipartForm.Value["tls_domain"][0]
		tlsKey = r.MultipartForm.Value["tls_domain_key"][0]
		if tls == "" || tlsKey == "" {
			return "", "", ErrFormDataParameterParsing
		}
		DOMAIN.Log.Info("Get tls from form data succeed!")
		return tls, tlsKey, nil
	}
	return "", "", ErrFormDataParameterParsing
}

func ValidTls(hostDomain, tlsCrt, tlsKey string) error {
	DOMAIN.Log.Info("Enter valid tls and tls key.")
	resultCrt, _ := Decode([]byte(tlsCrt))
	if resultCrt.Type != "CERTIFICATE" {
		return ErrInvalidTlsPem
	}
	pub, err := x509.ParseCertificate(resultCrt.Bytes)
	if err != nil {
		return ErrInvalidTlsPem
	}
	dnsName := pub.DNSNames
	isHostMatched := false
	for _, host := range dnsName {
		if hostDomain == host {
			isHostMatched = true
			break
		}
		ht := strings.HasPrefix(host, "*")
		if ht {
			hd := strings.Split(hostDomain, ".")
			hd[0] = "*"
			hdf := strings.Join(hd, ".")
			if hdf == host {
				isHostMatched = true
				break
			}
		}
	}
	DOMAIN.Log.Info("Got a ", pub, "with dns name :", dnsName, hostDomain)
	if !isHostMatched {
		return ErrInvalidTlsPem
	}
	_, err = tls.X509KeyPair([]byte(tlsCrt), []byte(tlsKey))
	if err != nil {
		return ErrInvalidTls
	}
	return nil
}
