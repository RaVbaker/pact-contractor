package s3

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/spf13/afero"

	"github.com/ravbaker/pact-contractor/internal/hooks"
	"github.com/ravbaker/pact-contractor/internal/parts"
	"github.com/ravbaker/pact-contractor/internal/paths"
	"github.com/ravbaker/pact-contractor/internal/speccontext"
)

const applicationJSON = "application/json"

func Upload(bucket, region, filesPath string, partContext parts.Context, ctx speccontext.GitContext) (err error) {
	var files []string
	var file afero.File

	client := NewClient(region)

	files, err = afero.Glob(fs, filesPath)
	if err != nil {
		return
	}

	if len(files) == 0 {
		return fmt.Errorf("no files found fo"+
			"r path: %q", filesPath)
	}
	log.Printf("For path %q detected files: %v", filesPath, files)
	for _, filename := range files {
		partContext = PrepareMergedFile(partContext, bucket, region, filename, ctx)
		file, err = fs.OpenFile(filename, os.O_RDONLY, 0400)
		if err != nil {
			return
		}
		path := paths.FilenameToPath(filename, partContext, ctx)
		_, err = upload(client, bucket, path, file, contextsToTagsMap(ctx, partContext))
		if err != nil {
			println(err.Error())
			return
		}
		err = hooks.Runner(path, filename, ctx, partContext)
		if err != nil {
			println("hook error", err.Error())
		}
	}
	return
}

func contextsToTagsMap(ctx speccontext.GitContext, partContext speccontext.PartsContext) map[string]*string {
	tags := make(map[string]*string)
	if len(ctx.Author) != 0 {
		tags["Author"] = aws.String(ctx.Author)
	}

	if len(ctx.CommitSHA) != 0 {
		tags["CommitSHA"] = aws.String(ctx.CommitSHA)
	}

	if len(ctx.Branch) != 0 {
		tags["Branch"] = aws.String(ctx.Branch)
	}

	if len(ctx.Origin) != 0 {
		tags["Origin"] = aws.String(ctx.Origin)
	}

	if partContext.Defined() && !partContext.Merged() {
		tags["Part"] = aws.String(partContext.Name())
	}
	return tags
}

// upload puts S3 object
// @TODO: upload only of content changed
func upload(client s3iface.S3API, bucket, path string, file afero.File, metadata map[string]*string) (*string, error) {
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

		Metadata:    metadata,
		ContentType: aws.String(applicationJSON),
	})
	if err != nil {
		// Print the error and exit.
		return nil, fmt.Errorf("unable to upload %q to %q, %w", path, bucket, err)
	}

	// Tagging Part name
	if partName, ok := metadata["Part"]; ok {
		versionId := *uploadedObject.VersionID
		err = Tag(client, bucket, path, versionId, map[string]*string{"Part": partName})
		if err != nil {
			log.Printf("Couldn't mark object %q#%q in %q with PartTag, error: %v", path, versionId, bucket, err)
		}
	}

	fmt.Printf("Successfully uploaded %q [version: %q] to %q\n", path, *uploadedObject.VersionID, bucket)
	return uploadedObject.VersionID, nil
}
