// Define a type Executable and one method to run it (handle stdout, err)
// Executable.Run(cmd string, arg[]string)
// Error handle:
// 			pathexec not found
//			Error starting Cmd
package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Executable struct {
	Name   string
	Params []string
}

func printCommand(c *exec.Cmd) {
	fmt.Printf("==> Executing: %s\n", strings.Join(c.Args, " "))
}

func printError(err error) {
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("==> Error: %s\n", err.Error()))
		os.Exit(1)
	}
}

func printOutput(o []byte) {
	fmt.Printf("==> Output: %s\n\n", string(o))
}

func runCommand(c *exec.Cmd) string {
	// Buffer output
	cOutput := &bytes.Buffer{}
	c.Stdout = cOutput

	// Run command
	err := c.Run()
	printError(err)

	// Display output and return it
	printCommand(c)

	o := cOutput.Bytes()
	printOutput(o)

	return string(o)
}

func (e Executable) Run() string {
	//Find executable path
	pathexec, err := exec.LookPath(e.Name)

	//Exit program if pathexec not found
	printError(err)

	//Execute cmd with all given arguments
	if n := len(e.Params); n > 0 {
		r := exec.Command(pathexec, e.Params...)
		o := runCommand(r)
		return string(o)
	}

	//Execute cmd without argument
	r := exec.Command(pathexec)
	o := runCommand(r)
	return string(o)
}

// A StopFunc is a function used to stop a resource that has previously been
// started and that is running in the background.
type StopFunc func() error

// StartDaemon starts a process and lets it run until the returned StopFunc is
// called.
func StartDaemon(name string, arg ...string) (StopFunc, error) {
	// If the process name is not an absolute path, looks for it.
	fullpath := name
	if !filepath.IsAbs(name) {
		var err error
		if fullpath, err = exec.LookPath(name); err != nil {
			return nil, err
		}
	}

	daemon := exec.Command(fullpath, arg...)
	daemon.Stdout = os.Stdout
	daemon.Stderr = os.Stderr
	if err := daemon.Start(); err != nil {
		return nil, fmt.Errorf("failed to start process %q: %s", name, err)
	}

	// Ok, the process has started, we now return a closure that will be
	// evaluated when the called of StartAsDaemon calls the returned function.
	stop := func() error {
		if err := daemon.Process.Kill(); err != nil {
			return fmt.Errorf("failed to kill process %q: %s", name, err)
		}
		// Wait that the process has actually been killed, but also that it's
		// finished to copy stuff to stdout/stderr.
		if err := daemon.Wait(); err != nil {
			return fmt.Errorf("exit error from process %q: %s", name, err)
		}
		return nil
	}
	return stop, nil
}

// PeriodicCall calls f at the given frequency and stop when stop is called.
func PeriodicCall(f func(), freq time.Duration) (StopFunc, error) {
	done := make(chan struct{})
	go func() {
		tick := time.NewTicker(freq)
		for {
			select {
			case <-tick.C:
				f()
			case <-done:
				log.Println("stopping periodic call")
			}
		}
	}()

	// Ok, the process has started, we now return a closure that will be
	// evaluated when the called of StartAsDaemon calls the returned function.
	stop := func() error {
		close(done)
		return nil
	}
	return stop, nil
}
