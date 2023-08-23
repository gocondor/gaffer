// Copyright © Harran Ali <harran.m@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var GenerateEventCmd = &cobra.Command{
	Use:   "gen:event [event-name]",
	Short: "Create an event",
	Long: `Helps you generate a boilderplate code for events
example:
gaffer gen:event my-event-name --job EventJobName

`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pwd, err := os.Getwd()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		eventName := args[0]
		eventJobName, err := cmd.Flags().GetString("job")
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		fPath := filepath.Join(pwd, "events/event-names.go")
		eventNamesFile, err := os.Open(fPath)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		eventNamesFileContent, err := io.ReadAll(eventNamesFile)
		eventNamesFile.Close()
		eventNamesFileContentStr := string(eventNamesFileContent)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		eventNamesFileContentStr = strings.TrimSuffix(eventNamesFileContentStr, "\n")
		eventNamesFileContentStr = eventNamesFileContentStr + "\n" + getEventNameStatement(prepareEventNameConst(eventName), eventName)
		os.WriteFile(fPath, []byte(eventNamesFileContentStr), 666)
		jfn := prepareJobFileName(eventJobName) + ".go"
		ffnp := filepath.Join(pwd, "events/jobs", jfn)
		jfs, err := os.Stat(ffnp)
		if err != nil && !os.IsNotExist(err) {
			fmt.Printf("problem creating the file: %v\n", ffnp)
			os.Exit(1)
		}

		if jfs != nil {
			fmt.Printf("file \"%v\" already exist\n", jfn)
			os.Exit(1)
		}
		jobFile, err := os.Create(ffnp)
		if err != err {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		jobFile.WriteString(prepareJobContent(eventJobName))
		jobFile.Close()
	},
}

func init() {
	rootCmd.AddCommand(GenerateEventCmd)
	GenerateEventCmd.Flags().StringP("job", "j", "", "the name of the job to be executed when the event is fired")
	GenerateEventCmd.MarkFlagRequired("job")
}
