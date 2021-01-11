package hooks

import (
	"fmt"
	"os"
	"os/exec"
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

type NoopHook struct{}

func (n NoopHook) Run(_ string) error {
	return nil
}

type LocalHook struct {
	Command string
}

func (l *LocalHook) Run(path string) error {
	cmd := strings.ReplaceAll(l.Command, "{path}", path)
	cmd = os.ExpandEnv(cmd)
	fmt.Printf("> %s\n", cmd)
	out, err := exec.Command("bash", "-c", cmd).CombinedOutput()
	if err != nil {
		return err
	}
	fmt.Println(string(out))
	return nil
}
