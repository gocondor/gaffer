// Copyright © Harran Ali <harran.m@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"github.com/spf13/cobra"
)

var GenerateEventJobCmd = &cobra.Command{
	Use:   "gen:eventjob",
	Short: "Create an event job",
	Long: `Helps you generate a boilderplate code for event jobs
example:
gaffer gen:eventjob EventJobName

`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO implement
	},
}

func init() {
	rootCmd.AddCommand(GenerateEventJobCmd)
}
