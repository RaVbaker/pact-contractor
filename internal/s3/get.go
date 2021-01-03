package s3

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func Get(bucket, region, path, versionID string) (fields map[string]*string, err error) {
	client := NewClient(region)
	fields, versionID, err = GetMetadata(client, bucket, path, versionID)
	if err != nil {
		return nil, err
	}
	fields["VersionID"] = &versionID

	tagDetails, err := client.GetObjectTagging(&s3.GetObjectTaggingInput{
		Bucket:    aws.String(bucket),
		Key:       aws.String(path),
		VersionId: optionalAWSString(versionID),
	})

	if err != nil {
		return fields, fmt.Errorf("cannot fetch tags for %q#%s from bucket %s, %w", path, versionID, bucket, err)
	}

	for _, tag := range tagDetails.TagSet {
		fields[*tag.Key] = tag.Value
	}

	return fields, nil
}

func GetMetadata(client *s3.S3, bucket, path string, versionID string) (map[string]*string, string, error) {
	object, err := client.HeadObject(&s3.HeadObjectInput{
		Bucket:    aws.String(bucket),
		Key:       aws.String(path),
		VersionId: optionalAWSString(versionID),
	})

	if err != nil {
		return nil, "", fmt.Errorf("cannot fetch metadata for %q#%s from bucket %s, %w", path, versionID, bucket, err)
	}

	return object.Metadata, *object.VersionId, err
}
