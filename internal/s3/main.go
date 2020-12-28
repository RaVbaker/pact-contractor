package s3

import (
	"github.com/spf13/afero"
)

const (
	DefaultSpecTag  = "main"
	defaultSpecName = "spec" // must match defaultFilesPath
)


var fs afero.Fs

func init() {
	fs = afero.NewOsFs()
}
