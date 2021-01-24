package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	// "syscall"
)

func run(input string) error {
	fmt.Fprintln(os.Stdout)

	commands := parseInput(input)
	var isError error = nil

CommandsLoop:
	for _, command := range commands {
		switch command[0] {
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
				return errors.New("")
			}

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

var lastDir string

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
