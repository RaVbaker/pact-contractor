package s3

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"
	
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/spf13/afero"
	
	"github.com/ravbaker/pact-contractor/internal/speccontext"
)

func Download(bucket, region, path, s3VersionID, gitBranch string, gitFlow bool) (err error) {
	paths := resolvePath(path, gitBranch, gitFlow)
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
		ok := download(bucket, region, path, s3VersionID, filename, file)
		if ok {
			break
		}
	}
	return
}

func resolvePath(path, gitBranch string, gitFlow bool) (paths []string) {
	if strings.Contains(path, speccontext.BranchSpecTag) {
		pattern := strings.Replace(path, speccontext.BranchSpecTag, "%s", -1)
		paths = gitBranchPaths(pattern, gitBranch, gitFlow)
	} else {
		paths = append(paths, path)
	}
	return
}

func gitBranchPaths(pattern, branchName string, gitFlow bool) (paths []string) {
	if branchName == "" {
		branchName = speccontext.CurrentBranchName()
	}
	if branchName != "" {
		paths = append(paths, fmt.Sprintf(pattern, branchName))
	}
	if gitFlow {
		paths = append(paths, fmt.Sprintf(pattern, speccontext.GitFlowDevelopBranch))
	}
	paths = append(paths, fmt.Sprintf(pattern, speccontext.DefaultSpecTag))
	return
}

func pathToFilename(path, spec string) string {
	dir, filename := filepath.Split(path)
	ext := filepath.Ext(filename)
	return filepath.Join(dir, spec+ext)
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
		Bucket: aws.String(bucket),
		Key:    aws.String(path),
		VersionId: optionalAWSString(s3VersionID),
	})
	if err != nil {
		log.Printf("failed to download file, %v", err)
		return false
	}
	if len(s3VersionID) != 0 {
		s3VersionID = fmt.Sprintf(" [version: %q]", s3VersionID)
	}
	fmt.Printf("Successfully downloaded %q%s from bucket %q to file %q, %d bytes\n", path, s3VersionID, bucket, filename, n)
	return true
}

func optionalAWSString(s string) *string {
	if len(s) == 0 {
		return nil
	}
	return aws.String(s)
}
