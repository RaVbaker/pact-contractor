package hooks

import (
	"path/filepath"
	"strings"

	"github.com/mitchellh/mapstructure"
)

type Hook struct {
	Type       string
	PathFilter string `mapstructure:"path_filter"`
	Spec       interface{}
}

type Spec interface {
	Run(path string) error
}

func (h *Hook) CanRun(path string) bool {
	if len(h.PathFilter) == 0 {
		return true
	}
	matched, err := filepath.Match(h.PathFilter, path)
	return err == nil && matched
}

func (h *Hook) Definition() Spec {
	switch strings.ToLower(h.Type) {
	case "local":
		var localSpec LocalHook
		mapstructure.Decode(h.Spec, &localSpec)
		return &localSpec
	}
	return &NoopHook{}
}
