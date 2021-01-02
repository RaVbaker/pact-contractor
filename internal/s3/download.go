package s3

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/spf13/afero"

	"github.com/ravbaker/pact-contractor/internal/paths"
)

func Download(bucket, region, pathsArg, s3VersionID, gitBranch string, gitFlow bool) (err error) {
	list := paths.Extract(pathsArg, s3VersionID)
	for path, version := range list {
		err = downloadPath(bucket, region, path, version, gitBranch, gitFlow)
		if err != nil {
			log.Printf("Download of \"%s#%s\", error %v", path, version, err)
		}
	}
	return
}

func downloadPath(bucket string, region string, path string, s3VersionID string, gitBranch string, gitFlow bool) (err error) {
	potentialPaths := paths.Resolve(path, gitBranch, gitFlow)
	af := afero.Afero{Fs: fs}
	var file afero.File
	var ok bool
	for _, potentialPath := range potentialPaths {
		filename := paths.PathToFilename(potentialPath)
		err = af.MkdirAll(filepath.Dir(filename), 0755)
		if err != nil {
			return fmt.Errorf("failed to create dir structure %q, %w", filepath.Dir(filename), err)
		}
		// Create a file to write the S3 Object contents to.
		file, err = af.Create(filename)
		if err != nil {
			return fmt.Errorf("failed to create file %q, %w", filename, err)
		}
		ok = download(bucket, region, potentialPath, s3VersionID, filename, file)
		if ok {
			break
		}
	}
	if !ok {
		err = fmt.Errorf("failed to download %q#%s from %s bucket into file", path, s3VersionID, bucket)

	}
	return
}

func download(bucket, region, path, s3VersionID, filename string, file afero.File) bool {
	// Initialize a session in us-west-2 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)

	// Setup the S3 Upload Manager. Also see the SDK doc for the Upload Manager
	// for more information on configuring part size, and concurrency.
	//
	// http://docs.aws.amazon.com/sdk-for-go/api/service/s3/s3manager/#NewUploader
	downloader := s3manager.NewDownloader(sess)

	// Write the contents of S3 Object to the file
	n, err := downloader.Download(file, &s3.GetObjectInput{
		Bucket:    aws.String(bucket),
		Key:       aws.String(path),
		VersionId: optionalAWSString(s3VersionID),
	})
	if err != nil {
		log.Printf("failed to download file %q from \"%s#%s\", %v", filename, path, s3VersionID, err)
		return false
	}
	if len(s3VersionID) != 0 {
		s3VersionID = fmt.Sprintf(" [version: %q]", s3VersionID)
	}
	fmt.Printf("Successfully downloaded %q%s from bucket %q to file %q, %d bytes\n", path, s3VersionID, bucket, filename, n)
	return true
}
