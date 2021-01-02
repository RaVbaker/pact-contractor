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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ravbaker/pact-contractor/internal/verification"
)

// checkCmd represents the canIDeploy command
var checkCmd = &cobra.Command{
	Use:     "can-i-deploy [path]",
	Aliases: []string{"get"},
	Short:   "Checks the verification status of the contract and displays details of a path, alias: get",
	Long: `Checks the verification status of the contract, when it is "success"
then return ExitCode 0, otherwise non-zero.

The command also prints details of the path, like metadata and tags.`,
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]

		return verification.CheckStatus(viper.GetString("bucket"), viper.GetString("region"), path, s3VersionID, gitBranchName, providerVersion, gitFlow)
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// checkCmd.PersistentFlags().String("foo", "", "A help for foo")

	checkCmd.Flags().StringVar(&s3VersionID, "version", "", "Provides AWS S3 Object VersionID for download")
	checkCmd.Flags().StringVar(&gitBranchName, "git-branch", "", "Overwrites git detected current branch name")
	checkCmd.Flags().BoolVar(&gitFlow, "gitflow", false, "Implies the Git Flow on matching branch detection(feature-branch,develop,main), if not enabled it uses GitHub Flow(feature-branch,main)")
	checkCmd.Flags().StringVarP(&providerVersion, "provider-version", "p", "", "Provider version/tag verified")
}
