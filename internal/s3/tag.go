package s3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

func Tag(client s3iface.S3API, bucket string, path, versionId string, tags map[string]*string) (err error) {
	var tagSet []*s3.Tag
	for name, value := range tags {
		tagSet = append(tagSet, &s3.Tag{Key: aws.String(name), Value: value})
	}
	_, err = client.PutObjectTagging(&s3.PutObjectTaggingInput{
		Bucket:    aws.String(bucket),
		Key:       aws.String(path),
		Tagging:   &s3.Tagging{TagSet: tagSet},
		VersionId: aws.String(versionId),
	})
	return
}
