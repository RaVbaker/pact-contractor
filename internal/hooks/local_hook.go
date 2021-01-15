package hooks

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

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
