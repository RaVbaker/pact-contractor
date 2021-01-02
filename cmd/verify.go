/*
Copyright © 2021 Rafal Piekarski <rafal.piekarski@hey.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// verifyCmd represents the verify command
var verifyCmd = &cobra.Command{
	Use:   "verify [path] --cmd [command]",
	Short: "Handy way to do pull && command && submit verification-status",
	Long: `Useful helper to execute common steps when contract is verified.

These steps are:
1. pulling the contract from S3 by [path]
2. running a command to verify the contract
3. submitting the verification status: success/failure

All flags possible for pull & submit are available for the verify command.
The path from first argument can be substituted in --cmd with "{path}" pattern,
which default value is for Ruby "bundle exec rake pact:verify:at[{path}]".
But feel free to set the "cmd" in config file for convenience.
Submitted status is detected from exitCode and 0 is interpreted as "success"
and any other as "failure".`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]
		pullCmd.Run(cmd, []string{path})

		cmdToRun := viper.GetString("cmd")
		const pathPattern = "{path}"
		cmdToRun = strings.ReplaceAll(cmdToRun, pathPattern, path)

		runCmd := exec.Command("bash", "-c", cmdToRun)
		out, err := runCmd.CombinedOutput()
		verificationStatus := "success"
		fmt.Printf("Executing command: `%s`\n\n%s\n", cmdToRun, out)
		if err != nil {
			verificationStatus = "failure"
			log.Printf("Command error: %v", err)
		}

		submitCmd.Run(cmd, []string{path, verificationStatus})
	},
}

var verificationCommand string

func init() {
	rootCmd.AddCommand(verifyCmd)
	verifyCmd.Flags().AddFlagSet(pullCmd.LocalFlags())
	verifyCmd.Flags().AddFlagSet(submitCmd.LocalFlags())

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// verifyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	verifyCmd.Flags().StringVar(&verificationCommand, "cmd", "bundle exec rake pact:verify:at[{path}]", "Command to execute during verification, {path} is replaced with provided [path].")
}
