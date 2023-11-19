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

var GenerateEventJobCmd = &cobra.Command{
	Use:   "gen:eventjob [JobName]",
	Short: "Create an event job",
	Long: `Helps you generate a boilderplate code for event jobs
example:
gaffer gen:eventjob EventJobName

`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pwd, err := os.Getwd()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		eventJobName := args[0]
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		jfn := camelCaseToSnake(eventJobName, "-") + ".go"
		ffnp := filepath.Join(pwd, "events/eventjobs", jfn)
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
		fmt.Println("event job generated successfully")
	},
}

func init() {
	rootCmd.AddCommand(GenerateEventJobCmd)
}
