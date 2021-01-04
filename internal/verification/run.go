package verification

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"

	"github.com/ravbaker/pact-contractor/internal/paths"
)

const pathPattern = "{path}"

func Run(cmdToRun, path, s3VersionID, gitBranchName string, cmd, pullCmd, submitCmd *cobra.Command) error {
	path, _ = paths.ForBranch(path, s3VersionID, gitBranchName)

	err := pullCmd.RunE(cmd, []string{path})
	if err != nil {
		return err
	}

	filename := paths.PathToFilename(path)
	cmdToRun = strings.ReplaceAll(cmdToRun, pathPattern, filename)

	fmt.Printf("Executing command: `%s`\n\n", cmdToRun)
	runCmd := exec.Command("bash", "-c", cmdToRun)
	runCmd.Stdout = os.Stdout
	runCmd.Stderr = os.Stderr
	err = runCmd.Run()

	verificationStatus := "success"
	if err != nil {
		verificationStatus = "failure"
		log.Printf("Command error: %v", err)
	}

	err = afero.Afero{Fs: fs}.Remove(filename)
	if err != nil {
		log.Printf("Couldn't remove downloaded file: %q, error: %v", filename, err)
	}

	err = submitCmd.RunE(cmd, []string{path, verificationStatus})

	return err
}
