package speccontext

import (
	"github.com/spf13/afero"
)

const (
	BranchSpecTag  = "{branch}"
	DefaultSpecTag = "main"
)

type Context struct {
	SpecTag            string
	Origin             string
	VerificationStatus bool
}

func NewContext(specTag, origin string) Context {
	return Context{SpecTag: specTag, Origin: origin}
}

var fs afero.Fs

func init() {
	fs = afero.NewOsFs()
}
