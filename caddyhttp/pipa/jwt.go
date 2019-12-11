package pipa

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
	StyleInfo
	jwt.StandardClaims
}

type StyleInfo struct {
	TimeStamp int64  `json:"time_stamp"`
	ProjectId string `json:"project_id"`
	Bucket    string `json:"bucket"`
}

func GetMethodFromJWT(r *http.Request, secretKey string) (claim *Claims, err error) {
	PIPA.Log.Println(10, "Enter get method from JWT")
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
	if claim.ProjectId == "" || claim.Bucket == "" {
		return claim, ErrJwtParameterParsing
	}
	PIPA.Log.Println(15, "Get the JWT parameters:", claim.ProjectId, claim.Bucket)
	return claim, nil
}
