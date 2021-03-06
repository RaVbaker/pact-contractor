/*
Copyright © 2020 Rafal Piekarski <rafal.piekarski@hey.com>

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

	"github.com/ravbaker/pact-contractor/internal/s3"
)

// pullCmd represents the pull command
var pullCmd = &cobra.Command{
	Use:   "pull [paths]",
	Short: "Pulls pact contracts from configured S3 bucket",
	Long: `Pulls pact contracts from the bucket by the path.

The path is a location of S3 object which base is the tag/git-branch name.
E.g. File stored under: "./pacts/provider/consumer/spec.json" can be fetched
for tag/branch: feature-xyz with command: "pull pacts/provider/consumer/spec/feature-xyz.json"

The current branch name can be substituted with "{branch}" so then the command would be:
"pull pacts/provider/consumer/spec/{branch}.json". When matching dynamic branch is not present
it uses GitHub Flow (feature-branch,main). Optional --gitflow is also possible to use the develop
as possible branch name.

The paths can be a list of paths separated by comma and with optional version definition after # sign. Like:
"paths/foo/bar/test.json#some-v3rsion-1D,paths/foo/baz/{branch}.json#oth3r-v3rsion-1D"
When paths are resolved with same values last definition is downloaded.
`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		paths := args[0]
		err := s3.Download(viper.GetString("bucket"), viper.GetString("region"), paths, s3VersionID, gitBranchName, gitFlow)
		return err
	},
}

var gitFlow bool
var gitBranchName, s3VersionID string

func init() {
	rootCmd.AddCommand(pullCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pullCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	pullCmd.Flags().StringVar(&s3VersionID, "version", "", "Provides AWS S3 Object VersionID for download")
	pullCmd.Flags().StringVar(&gitBranchName, "git-branch", "", "Overwrites git detected current branch name")
	pullCmd.Flags().BoolVar(&gitFlow, "gitflow", false, "Implies the Git Flow on matching branch detection(feature-branch,develop,main), if not enabled it uses GitHub Flow(feature-branch,main)")
}
