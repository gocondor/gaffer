// Copyright © Harran Ali <harran.m@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"github.com/spf13/cobra"
)

var GenerateHandlerCmd = &cobra.Command{
	Use:   "gen:handler",
	Short: "Create a handler function",
	Long: `Helps you generate a boilderplate code for a handler funtion
example
gaffer gen:handler ListUsers --file users.go

`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO implement
	},
}

func init() {
	rootCmd.AddCommand(GenerateHandlerCmd)
}
