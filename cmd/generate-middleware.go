// Copyright © Harran Ali <harran.m@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"github.com/spf13/cobra"
)

var GenerateMiddlewareCmd = &cobra.Command{
	Use:   "gen:middleware",
	Short: "Create a middleware",
	Long: `Helps you generate a boilderplate code for middlewares
example:
gaffer gen:middleware MyMiddleware

`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO implement
	},
}

func init() {
	rootCmd.AddCommand(GenerateMiddlewareCmd)
}
