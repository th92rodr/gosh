package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	// "syscall"
	"strings"
)

var lastDir string
var lastExitCode int

func (t *terminal) run(input string) error {
	fmt.Fprintln(os.Stdout)

	commands := parseInput(input)
	var isError error = nil

CommandsLoop:
	for _, command := range commands {
		switch command[0] {
		case backgroundOperator:
			goToGo := make(chan bool)
			go t.executeInBackground(command[1:], goToGo)
			<-goToGo

		case andOperator:
			if isError != nil {
				break CommandsLoop
			}

		case orOperator:
			if isError == nil {
				break CommandsLoop
			}

		case "exit":
			return errors.New("exit")

		case "cd":
			isError = cd(command[1])

		case "echo":
			if command[1] == "$?" {
				// Print the exit code of the last executed process.
				fmt.Fprintln(os.Stdout, lastExitCode)
			} else {
				isError = execute(command)
			}

		case semiColonOperator:
		case "":	// handle empty commands

		default:
			isError = execute(command)
		}
	}

	return nil
}

// Execute command in other process.
func execute(command []string) error {
	if binary, err := exec.LookPath(command[0]); err == nil {

		attr := new(os.ProcAttr)
		attr.Dir, _ = os.Getwd()
		attr.Env = os.Environ()
		attr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}

		if process, err := os.StartProcess(binary, command, attr); err == nil {
			processState, _ := process.Wait()

			// Get the process exit code,
			// in case of anything rather than 0, it means something unexpected happened.
			if exitCode := processState.ExitCode(); exitCode != 0 {
				lastExitCode = exitCode
				return errors.New("")
			}
			lastExitCode = 0

		} else {
			fmt.Fprintln(os.Stderr, err)
			return err
		}

	} else {
		fmt.Fprintln(os.Stderr, err)
		return err
	}

	return nil
}

func (t *terminal) executeInBackground(command []string, goToGo chan<- bool) {
	if binary, err := exec.LookPath(command[0]); err == nil {

		attr := new(os.ProcAttr)
		attr.Dir, _ = os.Getwd()
		attr.Env = os.Environ()
		attr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}

		if process, err := os.StartProcess(binary, command, attr); err == nil {
			fmt.Fprintln(os.Stdout, "[1]\t", process.Pid, "\t", strings.Join(command, " "))
			goToGo <- true
			process.Wait()
			fmt.Fprintln(os.Stdout, "\n[1]\t", process.Pid, "\t", strings.Join(command, " "), "\tDone")
			t.refresh()

		} else {
			fmt.Fprintln(os.Stderr, err)
			goToGo <- true
		}

	} else {
		fmt.Fprintln(os.Stderr, err)
		goToGo <- true
	}
}

// Execute command in the same process.
// func execute(command []string) {
// 	if binary, err := exec.LookPath(command[0]); err == nil {
// 		env := os.Environ()

// 		// Replaces the current process with the one invoked.
// 		if err = syscall.Exec(binary, command, env); err != nil {
// 			fmt.Fprintln(os.Stderr, err)
// 		}
// 	} else {
// 		fmt.Fprintln(os.Stderr, err)
// 	}
// }

// Change directory.
func cd(path string) error {
	currentDir, _ := os.Getwd()

	// If the informed path is a dash ("-") return to the last directory.
	if path == "-" && lastDir != "" {
		path = lastDir
		fmt.Fprintln(os.Stdout, path)
	}

	if err := os.Chdir(path); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}

	lastDir = currentDir

	return nil
}

func exit() {
	fmt.Fprintln(os.Stdout)
	os.Exit(0)
}
