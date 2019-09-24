package customdomaintest

import (
	. "github.com/journeymidnight/yig-front-caddy/test/lib"
	"net/http"
	"testing"
)

func Test_CustomDomain(t *testing.T) {
	sc := NewS3()
	defer sc.CleanEnv()
	sc.CleanEnv()
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

	domianHost := "www.wojiushixiangshiyishi.xyz"

	status, err := NewCustomDomain(domianHost)
	if status != http.StatusOK || err != nil {
		t.Fatal("Setting user-defined domain name failed:", status)
	}

	status, err = GetCustomDomain(domianHost)
	if status != http.StatusOK || err != nil {
		t.Fatal("Getting user-defined domain name failed:", status)
	}

	status, err = CustomDomainAccess(domianHost)
	if status != http.StatusOK || err != nil {
		t.Fatal("Custom domain access bucket resource failed:", status)
	}

	status, err = DelCustomDomain(domianHost)
	if status != http.StatusOK || err != nil {
		t.Fatal("Striking out user-defined domain name failed:", status)
	}

	sc.DeleteObject(TEST_BUCKET, TEST_KEY)
	sc.DeleteBucket(TEST_BUCKET)
}
