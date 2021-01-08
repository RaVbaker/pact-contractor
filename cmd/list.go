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
	"github.com/ravbaker/pact-contractor/internal/speccontext"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list [path]",
	Short: "List contracts for a given path",
	Long:  `List contracts for a given path (in a valid Glob format), sorted descending by creation date`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			args = append(args, defaultFilesPath)
		}

		return s3.List(viper.GetString("bucket"), viper.GetString("region"), args[0], specTag)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	listCmd.Flags().StringVarP(&specTag, "tag", "t", speccontext.BranchSpecTag, "Provides the tag under which is checked, file extension is ignored")
}
