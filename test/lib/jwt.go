package lib

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

func GetJwt(domainHost string) (string, error) {
	mySigningKey := []byte("hehehehe")
	type MyCustomClaims struct {
		TimeStamp    int64  `json:"time_stamp"`
		ProjectId    string `json:"project_id"`
		DomainHost   string `json:"host_domain"`
		Bucket       string `json:"bucket"`
		BucketDomain string `json:"bucket_domain"`
		jwt.StandardClaims
	}
	timeStamp := time.Now().Unix()
	bucketDomain := TEST_BUCKET + "." + End_Point
	// Create the Claims
	claims := MyCustomClaims{
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
