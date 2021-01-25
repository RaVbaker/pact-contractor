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
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ravbaker/pact-contractor/internal/hooks"
	"github.com/ravbaker/pact-contractor/internal/parts"
	"github.com/ravbaker/pact-contractor/internal/speccontext"
)

// testHookCmd represents the testHookCmd command
var testHookCmd = &cobra.Command{
	Use:   "test_hook [path]",
	Short: "When you want to debug hook execution use this command",
	Long:  `Can be used to test hook run`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]
		bucket := viper.GetString("bucket")
		fmt.Printf("Tested path: %q & bucket: %q / region: %q\n", path, bucket, regionName)
		fmt.Println("Lists AWS environment variables found in app config:")
		for _, key := range viper.AllKeys() {
			if strings.HasPrefix(key, "aws_") {
				safeValue := viper.GetString(key)
				if strings.Contains(key, "secret") && len(safeValue) > 0 {
					safeValue = "<SECRET>"
				}
				fmt.Printf("%q = %q\n", key, safeValue)
			}
		}

		fmt.Println("Lists AWS environment variables found globally:")
		for _, keyAndValue := range os.Environ() {
			if strings.HasPrefix(keyAndValue, "AWS_") {
				if strings.Contains(keyAndValue, "SECRET") {
					keyAndValue = strings.Split(keyAndValue, "=")[0] + "=<SECRET>"
				}
				fmt.Printf("%q\n", keyAndValue)
			}
		}

		if !hooks.Defined() {
			return fmt.Errorf("no hooks defined in config file")
		}

		fmt.Println("Trying to run hooks:")
		partsScope := parts.NewScope(part, numOfParts)
		ctx := speccontext.NewGitContext(specTag, contextOrigin, gitAuthor, gitBranch, gitCommitSHA)
		err := hooks.Runner(path, "test_filename", ctx, partsScope)
		if err != nil {
			return err
		}
		fmt.Println("Successfully pushed all hooks.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(testHookCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testHookCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// testHookCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	pushCmdFlagsSetup(testHookCmd)
}
