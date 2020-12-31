package verification

import (
	"fmt"
	"log"
	"os"
	
	"github.com/ravbaker/pact-contractor/internal/paths"
	"github.com/ravbaker/pact-contractor/internal/s3"
)

func PublishVerification(bucket, region, pathsArg, status, s3VersionID, providerVersion string)  {
	list := paths.Extract(pathsArg, s3VersionID)
	tags := tagsList(status, providerVersion)
	
	var err error
	for path, version := range list {
		err = s3.Tag(s3.NewClient(region), bucket, path, version, tags)
		if err != nil {
			log.Printf("Couldn't submit verification status %q for %q, error: %v", status, path, err)
		}
	}
	if err != nil {
		os.Exit(-1)
		return
	}
	fmt.Printf("Successfully marked as %q all paths: %q in bucket %s\n", status, pathsArg, bucket)
}

func tagsList(status string, providerVersion string) (list map[string]*string) {
	list = make(map[string]*string)
	list["Pact Verification"] =  &status
	if len(providerVersion) > 0 {
		list["Provider Version"] =  &providerVersion
	}
	return
}
