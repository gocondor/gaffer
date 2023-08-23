// Copyright © Harran Ali <harran.m@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT-style
// license that can be found in the LICENSE file.
package cmd

import (
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "gaffer",
	Version: "v1.5.1",
	Short:   "Gaffer is GoCondor's cli tool",
	Long:    `Gaffer is GoCondor's cli tool for creating new projects and performing other tasks`,

	// Run: func(cmd *cobra.Command, args []string) {

	// },
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".gaffer" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".gaffer")
	}

	viper.AutomaticEnv() // read in environment variables that match
}
