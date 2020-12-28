package s3

import (
	"fmt"
	"os"
	"path/filepath"
	
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/spf13/afero"
)

func Upload(bucket, region, specTag, filesPath string) (err error) {
	var files []string
	var file afero.File
	files, err = afero.Glob(fs, filesPath)
	for _, filename := range files {
		file, err = afero.Afero{Fs: fs}.OpenFile(filename, os.O_RDONLY, 0400)
		if err != nil {
			return
		}
		path := filenameToPath(filename, specTag)
		upload(bucket, region, path, file)
	}
	return
}

func filenameToPath(filename, tag string) string {
	dir, file := filepath.Split(filename)
	ext := filepath.Ext(file)
	return filepath.Join(dir, tag+ext)
}

func upload(bucket, region, path string, file afero.File) {
	// Initialize a session in us-west-2 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	
	// Setup the S3 Upload Manager. Also see the SDK doc for the Upload Manager
	// for more information on configuring part size, and concurrency.
	//
	// http://docs.aws.amazon.com/sdk-for-go/api/service/s3/s3manager/#NewUploader
	uploader := s3manager.NewUploader(sess)
	
	// Upload the file's body to S3 bucket as an object with the key being the
	// same as the filename.
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		
		// Can also use the `filepath` standard library package to modify the
		// filename as need for an S3 object key. Such as turning absolute path
		// to a relative path.
		Key: aws.String(path),
		
		// The file to be uploaded. io.ReadSeeker is preferred as the Uploader
		// will be able to optimize memory when uploading large content. io.Reader
		// is supported, but will require buffering of the reader's bytes for
		// each part.
		Body: file,
	})
	if err != nil {
		// Print the error and exit.
		panic(fmt.Sprintf("Unable to upload %q to %q, %v", path, bucket, err))
	}
	
	fmt.Printf("Successfully uploaded %q to %q\n", path, bucket)
}
