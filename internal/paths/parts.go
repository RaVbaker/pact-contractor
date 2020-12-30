package paths

import (
	"path/filepath"
	
	"github.com/ravbaker/pact-contractor/internal/speccontext"
)

func PathForPart(dir string, scope speccontext.PartsContext, ctx speccontext.GitContext) string {
	if scope.Merged() {
		return dir
	}
	return filepath.Join(dir, ctx.Branch, ctx.CommitSHA, scope.Name())
}

func CleanupPathForParts(filename string, ctx speccontext.GitContext) string {
	return filepath.Join(filepath.Dir(filename), ctx.Branch)
}