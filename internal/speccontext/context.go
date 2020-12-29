package speccontext

import (
	"github.com/spf13/afero"
)

const (
	BranchSpecTag = "{branch}"
	DefaultSpecTag  = "main"
)

type Context struct {
	SpecTag string
	VerificationStatus bool
}

func NewContext(specTag string) Context {
	return Context{SpecTag: specTag}
}

var fs afero.Fs

func init() {
	fs = afero.NewOsFs()
}
