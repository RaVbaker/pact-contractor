package parts

import (
	"encoding/json"
	"fmt"
	
	"github.com/ieee0824/go-deepmerge"
	"github.com/spf13/afero"
	
	"github.com/ravbaker/pact-contractor/internal/paths"
	"github.com/ravbaker/pact-contractor/internal/speccontext"
)

func Merge(filename string, partPaths []string, ctx speccontext.GitContext) error {
	currentPartContent, err := afero.ReadFile(fs, filename)
	if err != nil {
		return err
	}
	
	var mergedJSON, otherPart interface{}
	err = json.Unmarshal(currentPartContent, &mergedJSON)
	if err != nil {
		return fmt.Errorf("cannot unmarshal current part %q, %w", filename, err)
	}
	
	var otherFilename string
	var otherPartContent, serializedJSON []byte
	
	for _, path := range partPaths {
		otherFilename = paths.PathToFilename(path)
		otherPartContent, err = afero.ReadFile(fs, otherFilename)
		if err != nil {
			return err
		}
		
		if len(otherPartContent) == 0 {
			return fmt.Errorf("missing file %q", otherFilename)
		}
		
		err = json.Unmarshal(otherPartContent, &otherPart)
		if err != nil {
			return fmt.Errorf("cannot unmarshal part %q, %w", path, err)
		}
		mergedJSON, err = deepmerge.Merge(mergedJSON, otherPart)
	}
	
	serializedJSON, err = json.Marshal(mergedJSON)
	if err != nil {
		return err
	}
	
	err = afero.WriteFile(fs, filename, serializedJSON, 0755)
	if err != nil {
		return fmt.Errorf("cannot write merged file %q, %w", filename, err)
	}
	return nil
}

