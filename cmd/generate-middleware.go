// Copyright © Harran Ali <harran.m@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var GenerateMiddlewareCmd = &cobra.Command{
	Use:   "gen:middleware [MiddlewareName]",
	Short: "Create a middleware",
	Long: `Helps you generate a boilderplate code for middlewares
example:
gaffer gen:middleware MyMiddleware

`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pwd, err := os.Getwd()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		middlewareName := args[0]
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		mfn := camelCaseToSnake(middlewareName) + ".go"
		mfnp := filepath.Join(pwd, "middlewares", mfn)
		jfs, err := os.Stat(mfnp)
		if err != nil && !os.IsNotExist(err) {
			fmt.Printf("problem creating the file: %v\n", mfnp)
			os.Exit(1)
		}

		if jfs != nil {
			fmt.Printf("file \"%v\" already exist\n", mfn)
			os.Exit(1)
		}
		mwFile, err := os.Create(mfnp)
		if err != err {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		mwFile.WriteString(prepareMiddlewareContent(middlewareName))
		mwFile.Close()
	},
}

func init() {
	rootCmd.AddCommand(GenerateMiddlewareCmd)
}
