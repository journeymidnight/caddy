package lib

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type CustomDomainClaims struct {
	TimeStamp    int64  `json:"time_stamp"`
	ProjectId    string `json:"project_id"`
	DomainHost   string `json:"host_domain"`
	Bucket       string `json:"bucket"`
	BucketDomain string `json:"bucket_domain"`
	jwt.StandardClaims
}

type PipaClaims struct {
	TimeStamp int64  `json:"time_stamp"`
	ProjectId string `json:"project_id"`
	Bucket    string `json:"bucket"`
	jwt.StandardClaims
}

func GetJwtForCustomDomain(domainHost string) (string, error) {
	mySigningKey := []byte(SecretKey)
	timeStamp := time.Now().Unix()
	bucketDomain := TEST_BUCKET + "." + End_Point
	// Create the Claims
	claims := CustomDomainClaims{
		timeStamp,
		Project_Id,
		domainHost,
		TEST_BUCKET,
		bucketDomain,
		jwt.StandardClaims{},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", err
	}
	outPut := "Bearer " + ss
	fmt.Println("JWT is:", outPut)
	return outPut, nil
}

func GetJwtForPipa() (string, error) {
	mySigningKey := []byte(SecretKey)
	timeStamp := time.Now().Unix()
	// Create the Claims
	claims := PipaClaims{
		timeStamp,
		Project_Id,
		TEST_BUCKET,
		jwt.StandardClaims{},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", err
	}
	outPut := "Bearer " + ss
	return outPut, nil
}
