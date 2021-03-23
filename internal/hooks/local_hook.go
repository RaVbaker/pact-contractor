package hooks

import (
	"fmt"
	"os/exec"
)

type LocalHook struct {
	Command string
}

func (l *LocalHook) Run(path string) error {
	cmd := templateString(path, l.Command)
	fmt.Printf("> %s\n", cmd)
	out, err := exec.Command("sh", "-c", cmd).CombinedOutput()
	if err != nil {
		return err
	}
	fmt.Println(string(out))
	return nil
}
