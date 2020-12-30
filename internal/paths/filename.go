package paths

import (
	"path/filepath"
	
	"github.com/ravbaker/pact-contractor/internal/speccontext"
)

func FilenameToPath(filename string, partsScope speccontext.PartsContext, ctx speccontext.GitContext) string {
	dir, file := filepath.Split(filename)
	ext := filepath.Ext(file)
	
	if partsScope.Defined() {
		dir = PathForPart(dir, partsScope, ctx)
	}
	return filepath.Join(dir, ctx.SpecTag+ext)
}


