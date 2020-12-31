package s3

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/ravbaker/pact-contractor/internal/speccontext"
)

func List(bucket, region, path, specTag string) error {
	prefix := preparePrefix(path)
	pathPattern := strings.Replace(path, ".", "/*.", 1)

	client := NewClient(region)

	var contracts []*s3.Object

	err := client.ListObjectsPages(&s3.ListObjectsInput{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
	},
		func(page *s3.ListObjectsOutput, lastPage bool) bool {
			contracts = append(contracts, page.Contents...)
			return len(page.Contents) == int(*page.MaxKeys)
		})

	if err != nil {
		return err
	}

	sort.Slice(contracts, func(i, j int) bool {
		return contracts[i].LastModified.After(*contracts[j].LastModified)
	})

	var matchedPattern, matched bool
	for _, contract := range contracts {
		matchedPattern, err = filepath.Match(pathPattern, *contract.Key)
		matched, err = filepath.Match(path, *contract.Key)
		if matchedPattern || matched {
			fmt.Printf("%s\t%s\n", contract.LastModified, *contract.Key)
		}
	}

	return nil
}

func preparePrefix(path string) (prefix string) {
	if strings.Contains(path, "*") {
		components := strings.SplitN(path, "*", 2)
		prefix = components[0]
	} else {
		ext := filepath.Ext(path)
		prefix = strings.TrimSuffix(path, ext)
	}
	prefix = strings.Replace(prefix, speccontext.BranchSpecTag, "", -1)
	prefix = strings.TrimRight(prefix, "/")
	return
}
