package s3

import (
	"log"
	"strings"
	
	"github.com/ravbaker/pact-contractor/internal/parts"
	"github.com/ravbaker/pact-contractor/internal/paths"
	"github.com/ravbaker/pact-contractor/internal/speccontext"
)

func PrepareMergedFile(partContext parts.Context, bucket, region, filename string, ctx speccontext.GitContext) parts.Context {
	partPaths := downloadParts(partContext, bucket, region, filename, ctx)
	if len(partPaths) == 0 {
		return partContext
	}
	
	err := parts.Merge(filename, partPaths, ctx)
	if err != nil {
		log.Printf("Couldn't merge parts, error: %v", err)
		return partContext
	}
	partContext.MarkAsMerged()
	
	// 	cleanup locally and on S3
	fs.RemoveAll(paths.CleanupPathForParts(filename, ctx))
	for _, path := range partPaths {
		Delete(bucket, region, path)
	}
	return partContext
}

func downloadParts(scope parts.Context, bucket, region, filename string, ctx speccontext.GitContext) (partPaths []string) {
	for part := 1; part <= scope.Total(); part++ {
		path := paths.FilenameToPath(filename, parts.NewScope(part, scope.Total()), ctx)
		if part != scope.Current() {
			partPaths = append(partPaths, path)
		}
	}
	
	if len(partPaths) == 0 {
		return
	}
	
	log.Printf("Performing merge of %d parts. File %q from disk, and remotely from S3: %v", scope.Total(), filename, partPaths)
	
	err := Download(bucket, region, strings.Join(partPaths, ","), "", "", false)
	if err != nil {
		log.Printf("Cannot yet perform merge of parts (%d/%d) due to S3 fetch error: %v", scope.Current(), scope.Total(), err)
		return nil
	}
	
	return
}
