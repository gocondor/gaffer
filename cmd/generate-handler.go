// Copyright © Harran Ali <harran.m@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var GenerateHandlerCmd = &cobra.Command{
	Use:   "gen:handler [HandlerName]",
	Short: "Create a handler",
	Long: `Helps you generate a boilderplate code for handlers

example:
gaffer gen:handler ListUsers --file users.go

`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pwd, err := os.Getwd()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		handlerName := args[0]
		handlersFileName, err := cmd.Flags().GetString("file")
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		handlersFileName = strings.ToLower(handlersFileName)
		var handlersFile *os.File
		handlersFilePath := filepath.Join(pwd, "handlers", handlersFileName)
		hfs, err := os.Stat(handlersFilePath)
		if err != nil && !os.IsNotExist(err) {
			fmt.Printf("problem reading the file: %v\n", handlersFilePath)
			os.Exit(1)
		}

		if hfs == nil {
			handlersFile, err = createHandlerFile(handlersFilePath)
		} else {
			handlersFile, err = os.OpenFile(handlersFilePath, os.O_RDWR|os.O_APPEND, 766)
		}
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		handlersFile.WriteString(prepareHandlerContent(handlerName))
		fmt.Println("handler generated successfully")
	},
}

func init() {
	rootCmd.AddCommand(GenerateHandlerCmd)
	GenerateHandlerCmd.Flags().StringP("file", "f", "", "the file name within 'handlers/' directory in which to put the handler, if its not there a new one will be created")
	GenerateHandlerCmd.MarkFlagRequired("file")
}
