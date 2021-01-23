package hooks

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mitchellh/mapstructure"

	"github.com/ravbaker/pact-contractor/internal/parts"
	"github.com/ravbaker/pact-contractor/internal/speccontext"
)

type Hook struct {
	Type         string
	PathFilter   string            `mapstructure:"path_filter"`
	MatchContext map[string]string `mapstructure:"match_context"`
	RequireEnv   []string          `mapstructure:"require_env"`
	Spec         interface{}
}

type Spec interface {
	Run(path string) error
}

// CanRun ensures the hook can be executed. It checks that the path matches what specified in path_filter
// checks that match_context rules are satisfied. If the match defined with "/" as prefix and suffix it
// matches it as a regexp. To execute hook for a part upload you have to specify the match_context["part"]
// require_env matches can be a full line, specified like: "CI=true" or only with a name of the variable "CI"
func (h *Hook) CanRun(path string, ctx speccontext.GitContext, partsCtx parts.Context) (matched bool) {
	if len(h.PathFilter) > 0 {
		var err error
		matched, err = filepath.Match(h.PathFilter, path)
		if err != nil || !matched {
			return
		}
	}

	for key, value := range h.MatchContext {
		var field string
		switch key {
		case "branch":
			field = ctx.Branch
		case "author":
			field = ctx.Author
		case "tag":
			field = ctx.SpecTag
		case "origin":
			field = ctx.Origin
		case "commitSHA":
			field = ctx.CommitSHA
		case "part":
			field = partsCtx.Name()
		}
		matched = compareMatch(value, field)
		if !matched {
			return
		}
	}
	_, partRequirement := h.MatchContext["part"]
	if partsCtx.Defined() && !partsCtx.Merged() && !partRequirement {
		return false
	}

	for _, envDefinition := range h.RequireEnv {
		fullLine := strings.Contains(envDefinition, "=")
		if fullLine {
			matched = envInclude(envDefinition)
		} else {
			_, matched = os.LookupEnv(envDefinition)
		}
		if !matched {
			return
		}
	}

	return true
}

func envInclude(definition string) (matched bool) {
	for _, envLine := range os.Environ() {
		matched = matched || envLine == definition
	}
	return
}

func compareMatch(value, contextValue string) bool {
	if strings.HasPrefix(value, "/") && strings.HasSuffix(value, "/") {
		pattern := value[1 : len(value)-1]
		matched, err := regexp.MatchString(pattern, contextValue)
		if err != nil {
			log.Fatalf("Cannot match pattern: %q to %q, %v", value, contextValue, err)
		}
		return matched
	}
	return value == contextValue
}

func (h *Hook) Definition() Spec {
	switch strings.ToLower(h.Type) {
	case "local":
		var localSpec LocalHook
		mapstructure.Decode(h.Spec, &localSpec)
		return &localSpec
	case "http":
		var httpSpec HttpHook
		mapstructure.Decode(h.Spec, &httpSpec)
		return &httpSpec
	}
	return &NoopHook{}
}

func templateString(path, template string) string {
	template = strings.ReplaceAll(template, "{path}", path)
	return os.ExpandEnv(template)
}
