package handle

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strings"
	"time"
)

const DAYTIME = 86400

var Claim Claims

type Claims struct {
	CustomerInfo
	jwt.StandardClaims
}

type CustomerInfo struct {
	TimeStamp    int64  `json:"time_stamp"`
	ProjectId    string `json:"project_id"`
	DomainHost   string `json:"host_domain"`
	Bucket       string `json:"bucket"`
	BucketDomain string `json:"bucket_domain"`
}

func GetMethodFromJWT(r *http.Request, secretKey string) (claim *Claims, status int, err error) {
	tokenString := r.Header.Get("Authorization")
	tokenStrings := strings.Split(tokenString, " ")
	token, err := jwt.ParseWithClaims(tokenStrings[1], &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err == nil {
		claim, ok := token.Claims.(*Claims)
		time := time.Now().Unix()
		lastTime := time + DAYTIME
		if ok && token.Valid {
			if claim.TimeStamp <= lastTime {
				return claim, http.StatusOK, nil

			} else {
				return claim, http.StatusForbidden, fmt.Errorf("JWT:Token has expired")
			}
		} else {
			return claim, http.StatusInternalServerError, fmt.Errorf("JWT:Parameter conversion error")
		}
	} else {
		return claim, http.StatusForbidden, err
	}
}
