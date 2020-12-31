package s3

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func Delete(bucket, region, path string) bool {
	client := NewClient(region)

	_, err := client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		log.Printf("failed to delete key %q from bucket %q, %v", path, bucket, err)
		return false
	}

	fmt.Printf("Successfully deleted %q from bucket %q\n", path, bucket)
	return true
}
