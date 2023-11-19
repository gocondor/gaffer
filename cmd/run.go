// Copyright © Harran Ali <harran.m@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"github.com/radovskyb/watcher"
	"github.com/spf13/cobra"
)

var pid int

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run:dev",
	Short: "Start the app in hot reloading mode",
	Long: `To Start the app in hot reloading mode for development run the following command:

gaffer run:dev

`,
	Run: func(cmd *cobra.Command, args []string) {
		pwd, _ := os.Getwd()
		fileChangeChan := make(chan bool, 1)
		startAppChan := make(chan bool, 1)
		pidChan := make(chan int, 1)
		termSigsChan := make(chan os.Signal, 1)

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
		go startServerJob(pidChan, startAppChan)
		go startRestartControllerJob(fileChangeChan, pidChan, startAppChan)
		signal.Notify(termSigsChan, syscall.SIGINT, syscall.SIGTERM)
		go func(termSigsChan chan os.Signal) {
			for {
				<-termSigsChan
				if pid == 0 {
					os.Exit(0)
				} else {
					pgid, err := syscall.Getpgid(pid)
					if err != nil {
						fmt.Println("error getting pgid: ", err)
					}
					err = syscall.Kill(-pgid, syscall.SIGKILL)
					if err != nil {
						fmt.Println("error stopping process: ", err)
					}
					os.Exit(0)
				}
			}

		}(termSigsChan)

		func() {
			if err := w.Start(time.Millisecond * 100); err != nil {
				log.Fatalln(err)
			}
		}()
	},
}

func startRestartControllerJob(fileChangeChan chan bool, pidChan chan int, startAppChan chan bool) {
	for {
		fileChanged := <-fileChangeChan
		if fileChanged {
			fmt.Println("Restarting...")
			pid := <-pidChan
			if pid != 0 {
				if runtime.GOOS == "windows" {
					killCmd := exec.Command("taskkill", "/T", "/F", "/PID", strconv.Itoa(pid))
					err := killCmd.Run()
					if err != nil {
						fmt.Println("error stopping the dev server", err.Error())
					}
				} else if runtime.GOOS == "darwin" {
					pgid, err := syscall.Getpgid(pid)
					if err != nil {
						fmt.Println("error getting pgid: ", err.Error())
					}
					err = syscall.Kill(-pgid, syscall.SIGKILL)
					if err != nil {
						fmt.Println("error stopping process: ", err.Error())
					}
				} else {
					fmt.Println("not implemented for os: ", runtime.GOOS)
				}
				go func() {
					startAppChan <- true
				}()
			} else {
				go func() {
					startAppChan <- true
				}()
			}
			go func() {
				fileChangeChan <- false
			}()
		}
	}
}

func startServerJob(pidChan chan int, startAppChan chan bool) {
	for {
		shouldStartApp := <-startAppChan
		if shouldStartApp {
			fmt.Println("\n\nBuilding...")
			err := compileApp()
			if err != nil {
				fmt.Println("error building: ", err.Error())
				go func() {
					pidChan <- 0
				}()
			} else {
				fmt.Println("Starting...")
				pwd, _ := os.Getwd()
				if runtime.GOOS == "darwin" {
					execFile := pwd + "/tmp/" + path.Base(pwd)
					binary, lookErr := exec.LookPath("/bin/sh")
					if lookErr != nil {
						panic(lookErr)
					}
					args := []string{"/bin/sh", "-c", execFile}
					execSpec := &syscall.ProcAttr{
						Dir:   pwd,
						Env:   os.Environ(),
						Files: []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()},
						Sys: &syscall.SysProcAttr{
							Setpgid: true,
						},
					}
					pid, _, err = syscall.StartProcess(binary, args, execSpec)
					if err != nil {
						fmt.Println("error starting process: ", err.Error())
						go func() {
							pidChan <- 0
						}()
					} else {
						go func() {
							pidChan <- pid
						}()
					}
				} else {
					fmt.Println("not implemented for os: ", runtime.GOOS)
				}
			}
		}
	}
}

func compileApp() error {
	pwd, _ := os.Getwd()
	var command *exec.Cmd
	command = exec.Command("/bin/sh", "-c", fmt.Sprintf("go build -o %v/tmp/", pwd))
	command.Env = os.Environ()
	command.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	command.Dir = pwd
	o, err := command.CombinedOutput()
	if string(o) != "" {
		fmt.Println(string(o))
	}
	if err != nil {
		return err
	}
	return nil
}

func init() {
	rootCmd.AddCommand(runCmd)
}
