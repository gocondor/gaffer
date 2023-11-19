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

var GenerateModelCmd = &cobra.Command{
	Use:   "gen:model [ModelName]",
	Short: "Create a database model",
	Long: `Helps you generate a boilderplate code for database model
example:
gaffer gen:model User

`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pwd, err := os.Getwd()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		modelName := args[0]
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		mfn := camelCaseToSnake(modelName, "-") + ".go"
		mfnp := filepath.Join(pwd, "models/", mfn)
		mfs, err := os.Stat(mfnp)
		if err != nil && !os.IsNotExist(err) {
			fmt.Printf("problem creating the file: %v\n", mfnp)
			os.Exit(1)
		}

		if mfs != nil {
			fmt.Printf("file \"%v\" already exist\n", mfn)
			os.Exit(1)
		}
		ModelFile, err := os.Create(mfnp)
		if err != err {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		tableName := camelCaseToSnake(modelName, "_")
		tableName = singleToPlural(tableName)
		ModelFile.WriteString(prepareModelContent(modelName, tableName))
		ModelFile.Close()
		fmt.Println("model generated successfully")
	},
}

func init() {
	rootCmd.AddCommand(GenerateModelCmd)
}
