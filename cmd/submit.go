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

	"github.com/ravbaker/pact-contractor/internal/verification"
)

// submitCmd represents the verification command
var submitCmd = &cobra.Command{
	Aliases: []string{"submit-verification"},
	Use:     "submit [paths] [verification-status]",
	Short:   "Stores verification status in S3 Object Tag",
	Long: `Stores verification status in AWS S3 path object Tag.

The Tag is called "Pact Verification" and contains the [status] value.
Rules for [paths] are same as for pull command, so they can contain the S3Object VersionId.
Optional fields about provider can be set with flags.`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		err := verification.PublishVerification(viper.GetString("bucket"), viper.GetString("region"), args[0], args[1], verifiedS3VersionID, providerVersion, providerContext)
		return err
	},
}

var providerVersion, providerContext, verifiedS3VersionID string

func init() {
	rootCmd.AddCommand(submitCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// submitCmd.PersistentFlags().String("foo", "", "A help for foo")

	submitCmd.Flags().StringVarP(&providerVersion, "provider-version", "p", "", "Provider version/tag stored")
	submitCmd.Flags().StringVarP(&providerContext, "provider-context", "o", "", "Provides optional provider context (e.g. Build identifier or URL) value stored in Object Tags")
	submitCmd.Flags().StringVar(&verifiedS3VersionID, "version", "", "Provides AWS S3 Object VersionID for verification")
}
