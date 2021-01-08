package verification

import (
	"fmt"
	"log"

	"github.com/ravbaker/pact-contractor/internal/paths"
	"github.com/ravbaker/pact-contractor/internal/s3"
)

func CheckStatus(bucket, region, pathsArg, versionID, gitBranchName, expectedProviderVersion string, gitFlow bool) (err error) {
	var path string
	var fields map[string]*string

	list := paths.Extract(pathsArg, versionID)
	for extractedPath, version := range list {
		potentialPaths := paths.Resolve(extractedPath, gitBranchName, gitFlow)
		for _, potentialPath := range potentialPaths {
			fields, err = s3.Get(bucket, region, potentialPath, version)
			if err != nil {
				log.Printf("Get of \"%s#%s\", error %v", potentialPath, version, err)
			} else {
				path = potentialPath
				versionID = version
				break
			}
		}
		break
	}
	fmt.Printf("Examinating path: %q, version ID: %q\n\n", path, versionID)

	for key, value := range fields {
		fmt.Printf("%s: %q\n", key, *value)
	}

	if err != nil {
		return fmt.Errorf("couldn't fetch path details, %w", err)
	}

	versionID = *fields["VersionID"]

	providerVersionField, ok := fields["Provider Version"]
	if len(expectedProviderVersion) > 0 && ok && *providerVersionField != expectedProviderVersion {
		return fmt.Errorf("provider version mismatch, field %q != %q", *providerVersionField, expectedProviderVersion)
	}

	verifiedStatus, ok := fields["Pact Verification"]

	if !ok {
		return fmt.Errorf("no verification status found for %q#%s", path, versionID)
	} else if *verifiedStatus != "success" {
		return fmt.Errorf("unsuccessful verification, current status: %q", *verifiedStatus)
	}
	return nil
}
