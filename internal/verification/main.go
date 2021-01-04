package verification

import (
	"github.com/spf13/afero"
)

var fs afero.Fs

func init() {
	fs = afero.NewOsFs()
}
