// Copyright © 2021 Harran Ali <harran.m@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT-style
// license that can be found in the LICENSE file.
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gincoat",
	Short: "Gincoat framework installer",
	Long: `Gincoat installer helps you create new Gincoat framework projects.
`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("To create a new Gincoat project run the following command:")
		fmt.Println("gincoat new [project-name] [project-repository]")
		fmt.Println("for example:")
		fmt.Println("gincoat new my-app github.com/my-organization/my-app")

	},
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

		// Search config in home directory with name ".installer" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".gincoat")
	}

	viper.AutomaticEnv() // read in environment variables that match
}