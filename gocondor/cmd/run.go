// Copyright 2021 Harran Ali <harran.m@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"time"

	"github.com/radovskyb/watcher"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run:dev",
	Short: "Start the development server",
	Long: `To Start the development server run the following command:

cli run:dev

`,
	Run: func(cmd *cobra.Command, args []string) {
		pwd, _ := os.Getwd()
		fmt.Println("Starting ...")
		fileChangeChan := make(chan bool)
		processChan := make(chan *os.Process)

		w := watcher.New()
		w.SetMaxEvents(1)
		w.IgnoreHiddenFiles(true)
		w.Ignore(
			pwd + "/logs/app.log",
		)
		w.FilterOps(watcher.Rename, watcher.Move, watcher.Create, watcher.Write)
		if err := w.AddRecursive(pwd); err != nil {
			log.Fatalln(err)
		}

		go func(fileChangeChan chan bool) {
			for {
				select {
				case <-w.Event:
					fileChangeChan <- true
				case err := <-w.Error:
					log.Fatalln(err)
				case <-w.Closed:
					return
				}
			}
		}(fileChangeChan)

		go startServer(fileChangeChan, processChan)
		go restartController(fileChangeChan, processChan)

		// Start the watching process - it'll check for changes every 100ms.
		func() {
			if err := w.Start(time.Millisecond * 100); err != nil {
				log.Fatalln(err)
			}
		}()

	},
}

func restartController(fileChangeChan chan bool, processChan chan *os.Process) {
	process := <-processChan
	<-fileChangeChan
	fmt.Println("Restarting...")
	//recived file change, kill the process and then restart it again
	// stop for windows
	if runtime.GOOS == "windows" {
		killCmd := exec.Command("taskkill", "/T", "/F", "/PID", strconv.Itoa(process.Pid))
		err := killCmd.Run()
		if err != nil {
			fmt.Println("error stoping the server [os is windows]")
		}
	} else {
		// stop for other os
		err := process.Kill()
		if err != nil {
			fmt.Printf("error stoping the dev server [os is %s]", runtime.GOOS)
		}
	}

	// start the server again
	go startServer(fileChangeChan, processChan)

	// restart the controller
	go restartController(fileChangeChan, processChan)
}

func startServer(fileChangeChan chan bool, processChan chan *os.Process) {
	pwd, _ := os.Getwd()
	var command *exec.Cmd
	command = exec.Command("go", "run", "main.go")
	command.Dir = pwd
	stdout, err := command.StdoutPipe()
	if err != nil {
		fmt.Println("error getting a pipe to stdout")
		panic(err)
	}

	command.Start()
	processChan <- command.Process

	oneByte := make([]byte, 100)
	for {
		n, err := stdout.Read(oneByte)
		if err != nil {
			break
		}
		fmt.Println(string(oneByte[:n]))

	}
	command.Wait()
}

func init() {
	//rootCmd.AddCommand(runCmd)
}
