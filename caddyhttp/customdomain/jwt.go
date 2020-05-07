package customdomain

import (
	"github.com/dgrijalva/jwt-go"
	. "github.com/journeymidnight/yig-front-caddy/caddyerrors"
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

func GetMethodFromJWT(r *http.Request, secretKey string) (claim *Claims, err error) {
	DOMAIN.Log.Info("Enter get method from JWT")
	tokenString := r.Header.Get("Authorization")
	tokenStrings := strings.Split(tokenString, " ")
	token, err := jwt.ParseWithClaims(tokenStrings[1], &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return claim, ErrAccessDenied
	}
	claim, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return claim, ErrInvalidJwtParams
	}
	time := time.Now().Unix()
	lastTime := time + MINITETIME
	firstTime := time - DAYTIME
	if claim.TimeStamp > lastTime || claim.TimeStamp < firstTime {
		return claim, ErrExpiredToken
	}
	DOMAIN.Log.Info("Get the JWT parameters:", claim.ProjectId, claim.DomainHost, claim.Bucket, claim.BucketDomain)
	return claim, nil
}
