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
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ravbaker/pact-contractor/internal/parts"
	"github.com/ravbaker/pact-contractor/internal/s3"
	"github.com/ravbaker/pact-contractor/internal/speccontext"
)

const defaultFilesPath = "pacts/*-*.json"

var specTag, contextOrigin, gitAuthor, gitBranch, gitCommitSHA string
var part, numOfParts int

var pushCmd = &cobra.Command{
	Use:   "push [path]",
	Short: "Push generated pact contracts to configured S3 bucket, (default path=\"" + defaultFilesPath + "\")",
	Long: `Push generated pact contracts, based on path to configured S3 bucket.

Default path="` + defaultFilesPath + `", but can be configured until it's in Glob format.

E.g. For path like: "pacts/*-*.json" which could represent "pacts/provider-consumer.json" scenario
and under a branch named: feature-xyz the path for stored S3 object would be:
"pacts/provider-consumer/feature-xyz.json". So the extension (.json) would stay, the filename (provider-consumer)
will remain and the tag/branch-name would be added at the very end. Other context details like: author, commitSHA
and context origin will be persisted in object Metadata (can be accessed with "get" command).

When you're generating your contract in parts, due to complex specs setup you can provide --parts-total and --part
flags and eventually when all parts will get pushed it will merge it and store under appropriate S3 object.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			args = append(args, defaultFilesPath)
		}
		partsScope := parts.NewScope(part, numOfParts)
		ctx := speccontext.NewGitContext(specTag, contextOrigin, gitAuthor, gitBranch, gitCommitSHA)
		path := args[0]
		err := s3.Upload(viper.GetString("bucket"), viper.GetString("region"), path, partsScope, ctx)
		if err != nil {
			log.Fatalf("%v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pushCmd.PersistentFlags().String("foo", "", "A help for foo")

	pushCmdFlagsSetup(pushCmd)
}

func pushCmdFlagsSetup(cmd *cobra.Command) {
	cmd.Flags().IntVar(&part, "part", 0, "When provided as non-zero indicates the part which was pushed")
	cmd.Flags().IntVar(&numOfParts, "parts-total", 0, "When provided as non-zero indicates how many parts should be submitted, when all then it merges contract into a single file")
	cmd.Flags().StringVarP(&specTag, "tag", "t", speccontext.BranchSpecTag, "Provides the tag under which the specification is stored, if '"+speccontext.BranchSpecTag+"' uses Git current branch name")
	cmd.Flags().StringVarP(&contextOrigin, "context", "o", "", "Provides optional context origin (e.g. Build identifier or URL) value stored with S3 Object metadata")
	cmd.Flags().StringVar(&gitAuthor, "git-author", "", "Provides the git commit author name")
	cmd.Flags().StringVar(&gitBranch, "git-branch", "", "Provides the git current branch name")
	cmd.Flags().StringVar(&gitCommitSHA, "git-commit-sha", "", "Provides the git commit SHA reference, if provided can be an origin of author/branch name")
}
