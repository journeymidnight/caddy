package host

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"net/http"
)

type Claims struct {
	CustomerInfo
	jwt.StandardClaims
}

type CustomerInfo struct {
	ProjectId  string `json:"project-id"`
	DomainHost string `json:"host"`
	Bucket     string `json:"bucket"`
}

func GetMethodFromJWT(w http.ResponseWriter, r *http.Request, SecretKey string) (claim *Claims) {
	token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SecretKey), nil
		})

	if err == nil {
		claim, ok := token.Claims.(*Claims)
		if ok && token.Valid {
			return claim
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
	return
}
