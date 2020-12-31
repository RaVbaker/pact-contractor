package paths

import (
	"path/filepath"
	"strings"
	
	"github.com/ravbaker/pact-contractor/internal/speccontext"
)

func FilenameToPath(filename string, partsScope speccontext.PartsContext, ctx speccontext.GitContext) string {
	dir, file := filepath.Split(filename)
	ext := filepath.Ext(file)
	dir = filepath.Join(dir, strings.TrimSuffix(file, ext))
	
	if partsScope.Defined() {
		dir = PathForPart(dir, partsScope, ctx)
	}
	return filepath.Join(dir, ctx.SpecTag+ext)
}


