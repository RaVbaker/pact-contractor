package s3

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func Get(bucket, region, path, versionID string) (fields map[string]*string, err error) {
	client := NewClient(region)
	var lastModified *time.Time
	fields, versionID, lastModified, err = GetMetadata(client, bucket, path, versionID)
	if err != nil {
		return nil, err
	}
	fields["VersionID"] = &versionID
	fields["Last Modified"] = aws.String(lastModified.Format(time.RFC3339))

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

func GetMetadata(client *s3.S3, bucket, path string, versionID string) (map[string]*string, string, *time.Time, error) {
	object, err := client.HeadObject(&s3.HeadObjectInput{
		Bucket:    aws.String(bucket),
		Key:       aws.String(path),
		VersionId: optionalAWSString(versionID),
	})

	if err != nil {
		return nil, "", nil, fmt.Errorf("cannot fetch metadata for %q#%s from bucket %s, %w", path, versionID, bucket, err)
	}

	return object.Metadata, *object.VersionId, object.LastModified, err
}
