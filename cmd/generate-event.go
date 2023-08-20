// Copyright © Harran Ali <harran.m@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"github.com/spf13/cobra"
)

var GenerateEventCmd = &cobra.Command{
	Use:   "gen:event",
	Short: "Create an event",
	Long: `Helps you generate a boilderplate code for events
example:
gaffer gen:event my-event-name --job EventJobName

`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO implement
	},
}

func init() {
	rootCmd.AddCommand(GenerateEventCmd)
}
