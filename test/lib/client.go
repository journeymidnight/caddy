package lib

import (
	"github.com/journeymidnight/aws-sdk-go/aws"
	"github.com/journeymidnight/aws-sdk-go/aws/credentials"
	"github.com/journeymidnight/aws-sdk-go/aws/session"
	"github.com/journeymidnight/aws-sdk-go/service/s3"
)

type S3Client struct {
	Client *s3.S3
}

const (
	TEST_BUCKET      = "mybucket"
	TEST_KEY         = "testput"
	TEST_KEY_SPECIAL = "testputspecial:!@$%^&*()_+=-;?><| "
	TEST_VALUE       = "valueput"
	End_Point        = "s3.test.com"
	Project_Id       = "hehehehe"
	FileDir          = "/etc/caddy/"
)

func NewS3() *S3Client {
	creds := credentials.NewStaticCredentials("hehehehe", "hehehehe", "")

	// By default make sure a region is specified
	s3client := s3.New(session.Must(session.NewSession(
		&aws.Config{
			Credentials: creds,
			DisableSSL:  aws.Bool(true),
			Endpoint:    aws.String(End_Point),
			Region:      aws.String("r"),
		},
	),
	),
	)
	return &S3Client{s3client}
}
