package s3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/spf13/afero"
)

const (
	DefaultSpecTag  = "main"
	defaultSpecName = "spec" // must match defaultFilesPath
)


var fs afero.Fs

func init() {
	fs = afero.NewOsFs()
}


func newClient(region string) *s3.S3 {
	// Initialize a session in us-west-2 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	sess:= session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region)},
	))
	
	svc := s3.New(sess)
	return svc
}