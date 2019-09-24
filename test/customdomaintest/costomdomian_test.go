package customdomaintest

import (
	. "github.com/journeymidnight/yig-front-caddy/test/lib"
	"net/http"
	"testing"
)

func Test_CustomDomain(t *testing.T) {
	sc := NewS3()
	err := sc.MakeBucket(TEST_BUCKET)
	if err != nil {
		t.Fatal("Make Bucket err:", err)
	}
	err = sc.PutObject(TEST_BUCKET, TEST_KEY, TEST_VALUE)
	if err != nil {
		t.Fatal("Put Object err:", err)
	}
	err = sc.PutObjectAcl(TEST_BUCKET, TEST_KEY, ObjectCannedACLPublicRead)
	if err != nil {
		t.Fatal("PutObjectAcl err:", err)
	}
	urlNew := "http://" + End_Point + "/?x-oss-action=NewCustomDomain"
	token, err := GetJwt(Domain_Host)
	status, err := SetRequest("PUT", urlNew, token)
	if status != http.StatusOK || err != nil {
		t.Fatal("Setting user-defined domain name failed:", status)
	}
	urlGet := "http://" + End_Point + "/?x-oss-action=GetCustomDomain"
	status, err = SetRequest("GET", urlGet, token)
	if status != http.StatusOK || err != nil {
		t.Fatal("Getting user-defined domain name failed:", status)
	}
	urlDomain := "http://" + Domain_Host + "/" + TEST_KEY
	t.Log(urlDomain)
	status, err = SetRequestWithDomain("GET", urlDomain, "")
	if status != http.StatusOK || err != nil {
		t.Fatal("Custom domain access bucket resource failed:", status)
	}
	urlDelete := "http://" + End_Point + "/?x-oss-action=DelCustomDomain"
	status, err = SetRequest("DELETE", urlDelete, token)
	if status != http.StatusOK || err != nil {
		t.Fatal("Striking out user-defined domain name failed:", status)
	}
	sc.DeleteObject(TEST_BUCKET, TEST_KEY)
	sc.DeleteBucket(TEST_BUCKET)
}
