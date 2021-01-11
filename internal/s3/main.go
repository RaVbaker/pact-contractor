package s3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

var fs afero.Fs

func init() {
	fs = afero.NewOsFs()
}

func NewClient(region string) *s3.S3 {
	// Initialize a session in us-west-2 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	sess := NewSession(region, false)
	config := ClientConfig(region, viper.GetString("aws_assume_role_arn"), sess)
	svc := s3.New(sess, config)
	return svc
}

func NewSession(region string, debug bool) *session.Session {
	config := &aws.Config{Region: optionalAWSString(region)}
	config = config.WithCredentialsChainVerboseErrors(debug)

	return session.Must(session.NewSession(config))
}

func ClientConfig(region, assumeRoleArn string, sess *session.Session) *aws.Config {
	config := &aws.Config{Region: optionalAWSString(region)}
	if len(assumeRoleArn) > 0 {
		credentials := stscreds.NewCredentials(sess, assumeRoleArn)
		return config.WithCredentials(credentials)
	}
	return config
}

func optionalAWSString(s string) *string {
	if len(s) == 0 {
		return nil
	}
	return aws.String(s)
}
