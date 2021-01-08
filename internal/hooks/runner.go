package hooks

import (
	"github.com/ravbaker/pact-contractor/internal/parts"
	"github.com/ravbaker/pact-contractor/internal/speccontext"
)

func Runner(path, filename string, ctx speccontext.GitContext, partsCtx parts.Context) error {
	for _, hook := range config.Hooks {
		if hook.CanRun(path) {
			err := hook.Definition().Run(path)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
