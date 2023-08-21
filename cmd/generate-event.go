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
	"unicode"

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
		fmt.Println(eventNamesFileContentStr)
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

func getEventNameStatement(constName string, eventName string) string {
	t := `const {constName} = "{eventName}"`
	r := strings.Replace(t, "{constName}", constName, 1)
	r = strings.Replace(r, "{eventName}", eventName, 1)
	return r
}

func prepareEventNameConst(eventName string) string {
	var res string
	words := strings.Split(eventName, "-")
	for k, v := range words {
		if k == 0 {
			res = strings.ToUpper(v) + "_"
		} else if k < (len(words) - 1) {
			res = res + strings.ToUpper(v) + "_"
		} else {
			res = res + strings.ToUpper(v)
		}
	}
	return res
}

func prepareJobFileName(name string) string {
	var res string
	namesB := []byte(name)
	for i, v := range namesB {
		if i == 0 {
			res = res + strings.ToLower(string(v))
		} else {
			if !unicode.IsUpper(rune(v)) {
				res = res + string(v)
			} else {
				res = res + "-"
				res = res + strings.ToLower(string(v))
			}
		}
	}
	return res
}

func prepareJobContent(jobName string) string {
	t := `package eventjobs

import (
	"github.com/gocondor/core"
)

var {JobName} core.EventJob = func(event *core.Event, c *core.Context) {
	// logic implementation goes here...
}
`
	res := strings.Replace(t, "{JobName}", jobName, 1)
	return res
}

func init() {
	rootCmd.AddCommand(GenerateEventCmd)
	GenerateEventCmd.Flags().StringP("job", "j", "", "the name of the job to be executed when the event is fired")
	GenerateEventCmd.MarkFlagRequired("job")
}
