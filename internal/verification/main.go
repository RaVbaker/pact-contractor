package verification

import (
	"fmt"
	
	"github.com/ravbaker/pact-contractor/internal/paths"
	"github.com/ravbaker/pact-contractor/internal/s3"
)

func PublishVerification(bucket, region, pathsArg, status, s3VersionID, providerVersion string) (err error) {
	list := paths.Extract(pathsArg, s3VersionID)
	tags := tagsList(status, providerVersion)
	
	for path, version := range list {
		err = s3.Tag(s3.NewClient(region), bucket, path, version, tags)
		if err != nil {
			return
		}
	}
	
	fmt.Printf("Successfully marked as %q all paths: %q in bucket %s\n", status, pathsArg, bucket)
	return nil
}

func tagsList(status string, providerVersion string) (list map[string]*string) {
	list = make(map[string]*string)
	list["Pact Verification"] =  &status
	if len(providerVersion) > 0 {
		list["Provider Version"] =  &providerVersion
	}
	return
}
