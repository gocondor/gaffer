// Copyright © Harran Ali <harran.m@gmail.com>. All rights reserved.
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
	"unicode/utf8"

	"github.com/briandowns/spinner"
	"github.com/c4milo/unpackit"
	"github.com/karrick/godirwalk"
	"github.com/spf13/cobra"
	"github.com/thanhpk/randstr"
)

type Config struct {
	ReleaseUrl         string   `json:"releaseUrl"`
	CliReleasedVersion string   `json:"cliReleasedVersion"`
	Paths              []string `json:"paths"`
}

type RepoMeta struct {
	TagName    string `json:"tag_name"`
	TarBallUrl string `json:"tarball_url"`
}

// Config file
const CONFIG_URL string = "https://raw.githubusercontent.com/gocondor/gaffer/main/config.json"

const REPO_URL string = "https://api.github.com/repos/gocondor/gocondor/releases/latest"

// Temporary file name
var tempName string

// Current verson of the installer
var version string = "v1.5.1"

// struct for creating new project command
type CmdNew struct{}

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new [project-name] [project-repository]",
	Short: "Create a new gocondor projects",
	Long: `Create new gocondor projects, 
	
	Example:
	gaffer new myapp github.com/my-organization/myapp
`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		cn := CmdNew{}

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

		var config Config
		cn.DownloadConfig(&http.Client{}, CONFIG_URL, &config)
		// Check for update
		if yes := cn.IsUpdatedRequired(config.CliReleasedVersion); yes {
			fmt.Println(`
				This version of gaffer is outdated!
				Please update by running the following commands:
				
				go install github.com/gocondor/gaffer@latest
				
			`)
			os.Exit(0)
		}

		repoMeta := FetchRepoMeta(REPO_URL)
		downloadUrl := strings.Replace(config.ReleaseUrl, "{name}", repoMeta.TagName, 1)

		// Download the gocondor release
		fmt.Println("Downloading gocondor ...")
		filePath := cn.DownloadGoCondor(&http.Client{}, downloadUrl, cn.GenerateTempName())
		//Unpack file
		fmt.Println("Unpacking ...")
		pwd, _ := os.Getwd()
		cn.Unpack(filePath, pwd)

		// Rename to the user's given project name
		os.Rename("./gocondor-"+removeFirstCHar(repoMeta.TagName), "./"+projectName) //first char is `v`
		// Remove the downloaded gocondor archive
		os.Remove(filePath)

		projectPath := pwd + "/" + projectName

		//remove .github folder
		os.RemoveAll(projectPath + "/.github")

		// Fix imports
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

// Download gocondor archive
func (cn *CmdNew) DownloadGoCondor(http *http.Client, url string, tempName string) string {
	tempFilePath := os.TempDir() + "/" + tempName
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("error downloading the GoCondor release")
		os.Exit(1)
	}
	defer response.Body.Close()

	file, err := os.Create(tempFilePath)
	if err != nil {
		fmt.Println("error creating temp file")
		os.Exit(1)
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		fmt.Println("error writing the GoCondor release to file")
		os.Exit(1)
	}

	return tempFilePath
}

// Download config
func (cn *CmdNew) DownloadConfig(http *http.Client, url string, conf *Config) *Config {
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("error downloading config")
		os.Exit(1)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("error reading config")
		os.Exit(1)
	}

	err = json.Unmarshal(body, &conf)
	if err != nil {
		fmt.Println("error unmarshaling config")
		os.Exit(1)
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

func fixImports(dirName string, projectRepo string, paths []string) {
	err := godirwalk.Walk(dirName, &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			if !de.IsDir() && strings.Contains(osPathname, ".go") {
				file, err := ioutil.ReadFile(osPathname)
				if err != nil {
					fmt.Printf("error reading %s", osPathname)
					os.Exit(1)
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
		os.Exit(1)
	}
	newContent := strings.Replace(string(file), paths[0], projectRepo, -1)
	err = ioutil.WriteFile(dirName+"/go.mod", []byte(newContent), 0)
	if err != nil {
		fmt.Println("error writing to go.mod file")
		os.Exit(1)
	}
}

// Unpack GoCondor
func (cn *CmdNew) Unpack(filePath string, destPath string) {
	// Open file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("error opening the downloaded file")
		os.Exit(1)
	}
	defer file.Close()

	// Unpack it
	_, err = unpackit.Unpack(file, destPath)
	if err != nil {
		fmt.Println("error unpacking the downloaded release")
		os.Exit(1)
	}
}

// Generate random name
func (cn *CmdNew) GenerateTempName() string {
	return "gocondor_temp_" + randstr.Hex(8) + ".tar.gz"
}

func FetchRepoMeta(url string) RepoMeta {
	var repoMeta RepoMeta

	// get the latest released version number
	res, err := http.Get(url)
	if err != nil {
		os.Exit(1)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		os.Exit(1)
	}
	json.Unmarshal(body, &repoMeta)

	return repoMeta
}

func removeFirstCHar(str string) string {
	_, i := utf8.DecodeRuneInString(str)

	return str[i:]
}
