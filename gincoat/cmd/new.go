// Copyright © 2021 Harran Ali <harran.m@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/c4milo/unpackit"
	"github.com/karrick/godirwalk"
	"github.com/spf13/cobra"
	"github.com/thanhpk/randstr"
)

type Release struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}
type Config struct {
	Releases                 map[string]Release `json:"releases"`
	InstallerReleasedVersion string             `json:"installerReleasedVersion"`
	Paths                    []string           `json:"paths"`
}

// Config file
const CONFIG_URL string = "https://raw.githubusercontent.com/gincoat/installer/master/config.json"

// Temporary file name
var tempName string

// Current verson of the installer
var version string = "v0.1-beta.4"

// struct for creating new project command
type CmdNew struct{}

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new [project-name] [project-repository]",
	Short: "Create a new Gincoat projects",
	Long: `Create new Gincoat projects, 
	
Example:
  gincoat new my-app github.com/my-organization/my-app
`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		cn := CmdNew{}
		var config Config

		// Extract the args
		projectName := args[0]
		projectRepo := args[1]

		// show the spinner
		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()

		// Check if a directory with the given name exist
		_, err := os.Stat(projectName)
		if !os.IsNotExist(err) {
			fmt.Println("\nA directory with the given projct name alerady exist!")
			os.Exit(0)
		}

		// Download the config from github
		fmt.Println("Preparing ...")
		cn.DownloadConfig(&http.Client{}, CONFIG_URL, &config)
		selectedRelease := config.Releases["latest"]

		// Check for update
		if yes := cn.IsUpdatedRequired(config.InstallerReleasedVersion); yes {
			cn.PrintUpdateRequiredMessage()
		}

		// Download the Gincoat release
		fmt.Println("Downloading Gincoat ...")
		filePath := cn.DownloadGincoat(&http.Client{}, config.Releases["latest"].Url, cn.GenerateTempName())

		//Unpack file
		fmt.Println("Unpacking ...")
		pwd, _ := os.Getwd()
		cn.Unpack(filePath, pwd)

		// Rename to the user's given project name
		os.Rename("./"+selectedRelease.Name, "./"+projectName)

		// Remove the downloaded Gincoat archive
		os.Remove(filePath)

		// Fix imports
		projectPath := pwd + "/" + projectName
		fixImports(projectPath, projectRepo, config.Paths)

		// Run go mod tidy
		command := exec.Command("go", "mod", "tidy")
		command.Dir = projectPath
		stdout, err := command.StdoutPipe()
		command.Start()

		oneByte := make([]byte, 100)
		num := 1
		for {
			_, err := stdout.Read(oneByte)

			if err != nil {
				fmt.Printf(err.Error())
				break
			}
			r := bufio.NewReader(stdout)
			line, _, _ := r.ReadLine()
			fmt.Println(string(line))
			num = num + 1

			if num > 3 {
				os.Exit(0)
			}
		}
		command.Wait()

		// Hide the spinner
		s.Stop()

		fmt.Println("done!")
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}

// Download Gincoat archive
func (cn *CmdNew) DownloadGincoat(http *http.Client, url string, tempName string) string {
	tempFilePath := os.TempDir() + "/" + tempName
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("error downloading the Gincoat release")
		panic(err)
	}
	defer response.Body.Close()

	file, err := os.Create(tempFilePath)
	if err != nil {
		fmt.Println("error creating temp file")
		panic(err)
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		fmt.Println("error writing the Gincoat release to file")
		panic(err)
	}

	return tempFilePath
}

// Download config
func (cn *CmdNew) DownloadConfig(http *http.Client, url string, conf *Config) *Config {
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("error downloading config")
		panic(err)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("error reading config")
		panic(err)
	}

	err = json.Unmarshal(body, &conf)
	if err != nil {
		fmt.Println("error unmarshaling config")
		panic(err)
	}

	return conf
}

// Check for updates
func (cn *CmdNew) IsUpdatedRequired(LatestReleasedVersion string) bool {
	if LatestReleasedVersion != version {
		return true
	}
	return false
}

// Print update required message
func (cn *CmdNew) PrintUpdateRequiredMessage() {
	fmt.Println(`
	This version of the Gincoat installer is outdated!
	Please update by running the following commands:
	
	go get github.com/gincoat/gincoatinstaller
	
			`)
	os.Exit(1)
}

func fixImports(dirName string, projectRepo string, paths []string) {
	err := godirwalk.Walk(dirName, &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			if !de.IsDir() && strings.Contains(osPathname, ".go") {
				file, err := ioutil.ReadFile(osPathname)
				if err != nil {
					fmt.Printf("error reading %s", osPathname)
					panic(err)
				}
				newContent := strings.Replace(string(file), paths[0], projectRepo, -1)
				ioutil.WriteFile(osPathname, []byte(newContent), 0)
			}

			return nil
		},
		Unsorted: true,
	})
	if err != nil {
		fmt.Println("error scaning updating import paths")
		fmt.Println(err)
	}

	file, err := ioutil.ReadFile(dirName + "/go.mod")
	if err != nil {
		fmt.Println("error reading go.mod file")
		panic(err)
	}
	newContent := strings.Replace(string(file), paths[0], projectRepo, -1)
	err = ioutil.WriteFile(dirName+"/go.mod", []byte(newContent), 0)
	if err != nil {
		fmt.Println("error writing to go.mod file")
		panic(err)
	}
}

// Unpack Gincoat
func (cn *CmdNew) Unpack(filePath string, destPath string) {
	// Open file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("error opening the downloaded file")
		panic(err)
	}
	defer file.Close()

	// Unpack it
	_, err = unpackit.Unpack(file, destPath)
	if err != nil {
		fmt.Println("error unpacking the downloaded release")
		panic(err)
	}
}

// Generate random name
func (cn *CmdNew) GenerateTempName() string {
	return "gincoat_temp_" + randstr.Hex(8) + ".tar.gz"
}
