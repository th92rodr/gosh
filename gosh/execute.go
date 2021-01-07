package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	// "syscall"
)

func run(command []string) error {
	fmt.Fprintln(os.Stdout)

	switch command[0] {
	case "exit":
		return errors.New("exit")
	case "cd":
		cd(command[1])
	case "":	// handle empty commands
	default:
		execute(command)
	}

	return nil
}

// execute command in other process
func execute(command []string) {
	if binary, err := exec.LookPath(command[0]); err == nil {

		attr := new(os.ProcAttr)
		attr.Dir, _ = os.Getwd()
		attr.Env = os.Environ()
		attr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}

		if process, err := os.StartProcess(binary, command, attr); err == nil {
			process.Wait()
		} else {
			fmt.Fprintln(os.Stderr, err)
		}

	} else {
		fmt.Fprintln(os.Stderr, err)
	}
}

// execute command in the same process
// func execute(command []string) {
// 	binary, err := exec.LookPath(command[0])
// 	if err != nil {
// 		fmt.Fprintln(os.Stderr, err)
// 		return
// 	}

// 	env := os.Environ()

// 	// replaces the current process with the one invoked
// 	if err = syscall.Exec(binary, command, env); err != nil {
// 		fmt.Fprintln(os.Stderr, err)
// 		return
// 	}
// }

var lastDir string

// Change directory
func cd(path string) {
	currentDir, _ := os.Getwd()

	// if the informed path is a dash ("-") return to the last directory
	if path == "-" && lastDir != "" {
		path = lastDir
		fmt.Fprintln(os.Stdout, path)
	}

	if err := os.Chdir(path); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	lastDir = currentDir
}

func exit() {
	fmt.Fprintln(os.Stdout)
	os.Exit(0)
}
