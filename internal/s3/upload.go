package s3

import (
	"fmt"
	"os"
	
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/spf13/afero"
	
	"github.com/ravbaker/pact-contractor/internal/parts"
	"github.com/ravbaker/pact-contractor/internal/paths"
	"github.com/ravbaker/pact-contractor/internal/speccontext"
)

func Upload(bucket, region, filesPath string, partsScope parts.Scope, ctx speccontext.GitContext) (err error) {
	var files []string
	var file afero.File
	
	client := newClient(region)
	
	files, err = afero.Glob(fs, filesPath)
	for _, filename := range files {
		partsScope = PrepareMergedFile(partsScope, bucket, region, filename, ctx)
		file, err = fs.OpenFile(filename, os.O_RDONLY, 0400)
		if err != nil {
			return
		}
		path := paths.FilenameToPath(filename, partsScope, ctx)
		// consider parts to be uploaded with EXPIRATION DATE or
		upload(client, bucket, path, file, gitContextToTagsMap(ctx))
	}
	return
}

func gitContextToTagsMap(ctx speccontext.GitContext) map[string]*string {
	tags := make(map[string]*string)
	if len(ctx.Author) != 0 {
		tags["Author"] = &ctx.Author
	}
	
	if len(ctx.CommitSHA) != 0 {
		tags["CommitSHA"] = &ctx.CommitSHA
	}
	
	if len(ctx.Branch) != 0 {
		tags["Branch"] = &ctx.Branch
	}
	return tags
}

// upload puts S3 object
// @TODO: upload only of content changed
// @TODO: upload and do json.deepmerge if same revision exists for the branch
func upload(client s3iface.S3API, bucket, path string, file afero.File, metadata map[string]*string) *string {
	// Setup the S3 Upload Manager. Also see the SDK doc for the Upload Manager
	// for more information on configuring part size, and concurrency.
	//
	// http://docs.aws.amazon.com/sdk-for-go/api/service/s3/s3manager/#NewUploader
	uploader := s3manager.NewUploaderWithClient(client)
	
	// Upload the file's body to S3 bucket as an object with the key being the
	// same as the filename.
	uploadedObject, err := uploader.Upload(&s3manager.UploadInput{
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
		
		Metadata: metadata,
	})
	if err != nil {
		// Print the error and exit.
		panic(fmt.Sprintf("Unable to upload %q to %q, %v", path, bucket, err))
	}
	
	
	fmt.Printf("Successfully uploaded %q [version: %q] to %q\n", path, *uploadedObject.VersionID, bucket)
	return uploadedObject.VersionID
}
