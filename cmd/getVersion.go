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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ravbaker/pact-contractor/internal/paths"
	"github.com/ravbaker/pact-contractor/internal/s3"
)

// getVersionCmd represents the getVersion command
var getVersionCmd = &cobra.Command{
	Use:   "get-version [path]",
	Short: "Returns object S3 VersionID",
	Long: `Returns object S3 VersionID for the [path]

Use --silent/-q mode to nto print any logs.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var path string
		path, s3VersionID = paths.ForBranch(args[0], s3VersionID, "")
		bucket := viper.GetString("bucket")
		regionName := viper.GetString("region")
		_, versionID, _, err := s3.GetMetadata(s3.NewClient(regionName), bucket, path, s3VersionID)
		if err != nil {
			log.Fatalf("Cannot get Object %q#%s from bucket %q versionID, error: %v", path, s3VersionID, bucket, err)
		}
		fmt.Println(versionID)
	},
}

func init() {
	rootCmd.AddCommand(getVersionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getVersionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getVersionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	getVersionCmd.Flags().StringVar(&s3VersionID, "version", "", "Provides AWS S3 Object VersionID to verify")
}
