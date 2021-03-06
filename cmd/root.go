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
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"

	"github.com/ravbaker/pact-contractor/internal/hooks"
)

var (
	cfgFile    string
	silentMode bool
	bucketName string
	regionName string
)

const Version = "1.0.0"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Version: Version,
	Use:     "pact-contractor",
	Short:   "Pact Broker alternative based on top of AWS S3 storage",
	Long: `PactContractor - Pact Broker alternative based on top of AWS S3 storage

	Allows to pull and push pact contracts to/from an AWS S3 bucket
	`,
	SilenceUsage:  true,
	SilenceErrors: false,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version string) {
	if len(version) > 0 {
		rootCmd.Version = version
	}
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./(spec/)?pacts/.pact-contractor.yaml and $HOME/.pact-contractor.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&silentMode, "silent", "q", false, "Prevent STDERR to be populated")
	rootCmd.MarkPersistentFlagFilename("config", "yaml", "yml")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.PersistentFlags().StringVarP(&bucketName, "bucket", "b", "", "AWS S3 Bucket name")
	rootCmd.MarkPersistentFlagRequired("bucket")
	rootCmd.PersistentFlags().StringVarP(&regionName, "region", "r", "", "AWS S3 Region name")
	rootCmd.PersistentFlags().String("aws_assume_role_arn", "", "AWS AssumeRole ARN")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if silentMode {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home or local pacts directory with name ".pact-contractor" (without extension).
		viper.AddConfigPath("spec/pacts")
		viper.AddConfigPath("pacts")
		viper.AddConfigPath(home)
		viper.SetConfigName(".pact-contractor")
	}

	viper.SetEnvPrefix("pact")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Println("Using config file:", viper.ConfigFileUsed())
		hooks.Parse()
	}
	presetRequiredFlags(rootCmd)
	presetRequiredFlags(pushCmd)
	presetRequiredFlags(submitCmd)
	presetRequiredFlags(verifyCmd)
	loadAWSCredentialsToEnv()
}

func loadAWSCredentialsToEnv() {
	// bind AWS env vars
	viper.BindEnv("AWS_PROFILE")
	viper.BindEnv("AWS_ACCESS_KEY")
	viper.BindEnv("AWS_ACCESS_KEY_ID")
	viper.BindEnv("AWS_SECRET_KEY")
	viper.BindEnv("AWS_SECRET_ACCESS_KEY")
	viper.BindEnv("AWS_SESSION_TOKEN")
	viper.BindEnv("AWS_CONFIG_FILE")
	viper.BindEnv("AWS_SHARED_CREDENTIALS_FILE")
	viper.BindEnv("AWS_ROLE_ARN")
	viper.BindEnv("AWS_CA_BUNDLE")
	viper.BindEnv("AWS_REGION")
	viper.BindEnv("AWS_DEFAULT_REGION")
	viper.BindEnv("AWS_SDK_LOAD_CONFIG")

	for _, key := range viper.AllKeys() {
		value := strings.TrimSpace(viper.GetString(key))
		if strings.HasPrefix(key, "aws_") && len(value) > 0 {
			envKey := strings.ToUpper(key)
			os.Setenv(envKey, value)
		}
	}
}

func presetRequiredFlags(cmd *cobra.Command) {
	viper.BindPFlags(cmd.Flags())
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if viper.IsSet(f.Name) && viper.GetString(f.Name) != "" {
			cmd.Flags().Set(f.Name, viper.GetString(f.Name))
		}
	})
}
