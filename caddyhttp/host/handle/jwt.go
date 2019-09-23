package handle

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strings"
	"time"
)

const DAYTIME = 86400
const MINITETIME = 900

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
	HOST.Log.Println(10, "Enter get method from JWT")
	tokenString := r.Header.Get("Authorization")
	tokenStrings := strings.Split(tokenString, " ")
	token, err := jwt.ParseWithClaims(tokenStrings[1], &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return claim, http.StatusForbidden, err
	}
	claim, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return claim, http.StatusInternalServerError, fmt.Errorf("JWT:Parameter conversion error")
	}
	time := time.Now().Unix()
	lastTime := time + MINITETIME
	firstTime := time - DAYTIME
	if claim.TimeStamp > lastTime || claim.TimeStamp < firstTime {
		return claim, http.StatusForbidden, fmt.Errorf("JWT:Token has expired")
	}
	HOST.Log.Println(15, "Get the JWT parameters:", claim.ProjectId, claim.DomainHost, claim.Bucket, claim.BucketDomain)
	return claim, http.StatusOK, nil
}
