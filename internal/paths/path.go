package paths

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	
	"github.com/ravbaker/pact-contractor/internal/speccontext"
)

const (
	DefaultSpecName = "spec" // must match defaultFilesPath
)


func Extract(paths string, versionID string) (out map[string]string) {
	out = make(map[string]string)
	pathList := strings.Split(paths, ",")
	for _, pathWithVersion := range pathList {
		splitPath := strings.SplitN(pathWithVersion+"#"+versionID, "#", 3)
		out[splitPath[0]] = splitPath[1]
	}
	return
}

func Resolve(path, gitBranch string, gitFlow bool) (paths []string) {
	if strings.Contains(path, speccontext.BranchSpecTag) {
		pattern := strings.Replace(path, speccontext.BranchSpecTag, "%s", -1)
		paths = gitBranchPaths(pattern, gitBranch, gitFlow)
	} else {
		paths = append(paths, path)
	}
	return
}

func gitBranchPaths(pattern, branchName string, gitFlow bool) (paths []string) {
	if branchName == "" {
		branchName = speccontext.CurrentBranchName()
	}
	if branchName != "" {
		paths = append(paths, fmt.Sprintf(pattern, branchName))
	}
	if gitFlow {
		paths = append(paths, fmt.Sprintf(pattern, speccontext.GitFlowDevelopBranch))
	}
	paths = append(paths, fmt.Sprintf(pattern, speccontext.DefaultSpecTag))
	return
}

func PathToFilename(path string) string {
	dir, tagFilename := filepath.Split(path)
	ext := filepath.Ext(tagFilename)
	return strings.TrimRight(dir, string(os.PathSeparator))+ext
}


