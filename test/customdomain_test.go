package test

import (
	. "github.com/journeymidnight/yig-front-caddy/test/lib"
	"net/http"
	"testing"
	"time"
)

const (
	DomainHost = "www.wojiushixiangshiyishi.xyz"
	TIMEWAIT = 10*time.Second
)

func Test_CustomDomain(t *testing.T) {
	sc := NewS3()
	defer sc.CleanEnv()
	defer DelCustomDomain(DomainHost)
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

	status, err := NewCustomDomain(DomainHost)
	if status != http.StatusOK || err != nil {
		t.Fatal("Setting user-defined domain name failed:", status)
	}

	status, err = GetCustomDomain(DomainHost)
	if status != http.StatusOK || err != nil {
		t.Fatal("Getting user-defined domain name failed:", status)
	}

	status, err = CustomDomainAccess(DomainHost)
	if status != http.StatusOK || err != nil {
		t.Fatal("Custom domain access bucket resource failed:", status)
	}

	status, err = DelCustomDomain(DomainHost)
	if status != http.StatusOK || err != nil {
		t.Fatal("Striking out user-defined domain name failed:", status)
	}

	sc.DeleteObject(TEST_BUCKET, TEST_KEY)
	sc.DeleteBucket(TEST_BUCKET)
}

func Test_CustomDomainWithTls(t *testing.T) {
	sc := NewS3()
	defer sc.CleanEnv()
	defer DelCustomDomain(DomainHost)
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

	status, err := NewCustomDomain(DomainHost)
	if status != http.StatusOK || err != nil {
		t.Fatal("Setting user-defined domain name failed:", status)
	}

	status, err = GetCustomDomain(DomainHost)
	if status != http.StatusOK || err != nil {
		t.Fatal("Getting user-defined domain name failed:", status)
	}

	err = PostTlsPem(DomainHost)
	if err != nil {
		t.Fatal("Custom domain Custom domain name add certificate:", status)
	}

	time.Sleep(TIMEWAIT)
	status, err = CustomDomainAccessWithcert(DomainHost)
	if status != http.StatusOK || err != nil {
		t.Fatal("Custom domain access bucket resource failed:", status)
	}

	status, err = DelTlsPem(DomainHost)
	if status != http.StatusOK || err != nil {
		t.Fatal("Striking out user-defined domain name tls pem failed:", status)
	}

	status, err = DelCustomDomain(DomainHost)
	if status != http.StatusOK || err != nil {
		t.Fatal("Striking out user-defined domain name failed:", status)
	}

	sc.DeleteObject(TEST_BUCKET, TEST_KEY)
	sc.DeleteBucket(TEST_BUCKET)
}
