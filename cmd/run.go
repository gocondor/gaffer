// Copyright 2021 Harran Ali <harran.m@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strconv"
	"syscall"
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
		fileChangeChan := make(chan bool, 5)
		startAppChan := make(chan bool, 5)
		stdoutChan := make(chan io.ReadCloser, 5)
		cmdChan := make(chan *exec.Cmd)
		fileChangeChan <- false
		startAppChan <- false
		w := watcher.New()
		w.SetMaxEvents(1)
		w.IgnoreHiddenFiles(true)
		w.Ignore(
			pwd+"/logs/app.log",
			pwd+"/tmp",
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
		go func() {
			startAppChan <- true
		}()
		go startServerJob(cmdChan, startAppChan, stdoutChan)
		go startRestartControllerJob(fileChangeChan, cmdChan, startAppChan, stdoutChan)

		func() {
			if err := w.Start(time.Millisecond * 100); err != nil {
				log.Fatalln(err)
			}
		}()
	},
}

func startRestartControllerJob(fileChangeChan chan bool, cmdChan chan *exec.Cmd, startAppChan chan bool, stdoutChan chan io.ReadCloser) {
	for {
		fileChanged := <-fileChangeChan
		if fileChanged {
			fmt.Println("Restarting...")
			startCmd := <-cmdChan
			if runtime.GOOS == "windows" {
				killCmd := exec.Command("taskkill", "/T", "/F", "/PID", strconv.Itoa(startCmd.Process.Pid))
				err := killCmd.Run()
				if err != nil {
					fmt.Printf("error stoping the dev server")
				}
			} else if runtime.GOOS == "darwin" {
				// fmt.Println("stop process state:", startCmd.ProcessState.String())
				err := startCmd.Process.Kill()
				if err != nil {
					fmt.Println("error stoping the dev server ")
				}
			}
			go func() {
				startAppChan <- true
			}()
			go func() {
				fileChangeChan <- false
			}()
			startCmdStdout := <-stdoutChan
			if startCmdStdout != nil {
				startCmdStdout.Close()
			}
		}
	}
}

func startServerJob(cmdChan chan *exec.Cmd, startAppChan chan bool, stdoutChan chan io.ReadCloser) {
	for {
		shouldStartApp := <-startAppChan
		if shouldStartApp {
			fmt.Println("\n\nBuilding...")
			compileApp()
			time.Sleep(time.Microsecond * 100)
			pwd, _ := os.Getwd()
			var command *exec.Cmd
			execFile := pwd + "/tmp/" + path.Base(pwd)
			fmt.Println("Starting...")
			command = exec.Command("/bin/sh", "-c", execFile)
			command.Env = os.Environ()
			command.Dir = pwd
			command.SysProcAttr = &syscall.SysProcAttr{
				Setpgid: true,
			}
			stdout, err := command.StdoutPipe()
			if err != nil {
				fmt.Println("error getting a pipe to stdout")
				panic(err)
			}
			err = command.Start()
			// fmt.Println("start process state:", command.ProcessState.String())
			if err != nil {
				fmt.Println("error starting the app ", err.Error())
				panic(err)
			}
			go func() {
				cmdChan <- command
			}()
			go func() {
				stdoutChan <- stdout
			}()
			oneByte := make([]byte, 100)
			for {
				n, err := stdout.Read(oneByte)
				if err != nil {
					break
				}
				fmt.Println(string(oneByte[:n]))
			}
			go func() {
				startAppChan <- false
			}()
			command.Wait()
		}
	}
}

func compileApp() {
	pwd, _ := os.Getwd()
	var command *exec.Cmd
	command = exec.Command("/bin/sh", "-c", fmt.Sprintf("go build -o %v/tmp/", pwd))
	command.Env = os.Environ()
	command.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	command.Dir = pwd
	stdout, err := command.StdoutPipe()
	defer stdout.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	err = command.Start()
	if err != nil {
		fmt.Println(err.Error())
	}
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
	rootCmd.AddCommand(runCmd)
}