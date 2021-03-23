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

func Run(cmdToRun, pathsArg, s3VersionID, gitBranchName string, cmd, pullCmd, submitCmd *cobra.Command) (err error) {
	var path string
	list := paths.Extract(pathsArg, s3VersionID)
	for pathNoVersionId, version := range list {
		potentialPaths := paths.Resolve(pathNoVersionId, gitBranchName, false)
		for _, potentialPath := range potentialPaths {
			pullCmd.Flag("version").Value.Set(version)
			err = pullCmd.RunE(cmd, []string{potentialPath})
			if err == nil {
				path = potentialPath
				s3VersionID = version
				break
			}
		}
		break
	}

	if err != nil && len(path) == 0 {
		return err
	}

	// path, _ = paths.ForBranch(pathsArg, s3VersionID, gitBranchName)

	filename := paths.PathToFilename(path)
	cmdToRun = strings.ReplaceAll(cmdToRun, pathPattern, filename)

	fmt.Printf("Executing command: `%s`\n\n", cmdToRun)
	runCmd := exec.Command("sh", "-c", cmdToRun)
	runCmd.Stdout = os.Stdout
	runCmd.Stderr = os.Stderr
	runErr := runCmd.Run()

	verificationStatus := "success"
	if runErr != nil {
		verificationStatus = "failure"
		log.Printf("Command error: %v", err)
	}

	err = afero.Afero{Fs: fs}.Remove(filename)
	if err != nil {
		log.Printf("Couldn't remove downloaded file: %q, error: %v", filename, err)
	}

	err = submitCmd.RunE(cmd, []string{path, verificationStatus})

	if err != nil {
		return err
	}

	return runErr
}
