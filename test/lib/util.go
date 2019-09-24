package lib

import (
	"net/url"
	"os"
	"strings"
)

func GenTestObjectUrl(sc *S3Client) string {
	return "http://" + *sc.Client.Config.Endpoint + string(os.PathSeparator) + TEST_BUCKET + string(os.PathSeparator) + TEST_KEY
}

func GenTestSpecialCharaterObjectUrl(sc *S3Client) string {
	urlchange := url.QueryEscape(TEST_KEY_SPECIAL)
	urlchange = strings.Replace(urlchange, "+", "%20", -1)
	return "http://" + *sc.Client.Config.Endpoint + string(os.PathSeparator) + TEST_BUCKET + string(os.PathSeparator) + urlchange
}

func (sc *S3Client) CleanEnv() {
	sc.DeleteObject(TEST_BUCKET, TEST_KEY)
	sc.DeleteBucket(TEST_BUCKET)
}

type AccessPolicyGroup struct {
	BucketPolicy string
	BucketACL    string
	ObjectACL    string
}

// Generate 5M part data
func GenMinimalPart() []byte {
	return RandBytes(5 << 20)
}
