package s3

import (
	"fmt"
	"path/filepath"
	
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/spf13/afero"
)

func Download(bucket, region, path string) (err error) {
	paths := []string{ path }
	af := afero.Afero{Fs: fs}
	var file afero.File
	for _, path := range paths {
		filename := pathToFilename(path, defaultSpecName)
		err = af.MkdirAll(filepath.Dir(filename), 0755)
		if err != nil {
			return fmt.Errorf("failed to create dir structure %q, %v", filepath.Dir(filename), err)
		}
		// Create a file to write the S3 Object contents to.
		file, err = af.Create(filename)
		if err != nil {
			return fmt.Errorf("failed to create file %q, %v", filename, err)
		}
		download(bucket, region, path, filename, file)
	}
	return
}

func pathToFilename(path, spec string) string {
	dir, filename := filepath.Split(path)
	ext := filepath.Ext(filename)
	return filepath.Join(dir, spec+ext)
}

func download(bucket, region, path, filename string, file afero.File) {
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
		Bucket: aws.String(bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		panic(fmt.Errorf("failed to download file, %v", err))
	}
	fmt.Printf("file downloaded, %d bytes\n", n)
}
