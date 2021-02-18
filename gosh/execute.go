package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

func (t *terminal) run(input string) error {
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

		case backgroundOperator:
			t.background(command)

		case "exit":
			return errors.New("exit")

		case "cd":
			isError = t.cdCommand(command)

		case "echo":
			isError = t.echo(command)

		case "fg":
			t.foreground()

		case semiColonOperator:
		case "":	// handle empty commands
		default:
			isError = t.execute(command)
		}
	}

	return nil
}

// Execute command in other process.
func (t *terminal) execute(command []string) error {
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
				t.lastExitCode = exitCode
				return errors.New("")
			}
			t.lastExitCode = 0

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

func (t *terminal) background(command []string) {
	if len(command) > 1 && (command[1] == "cd" || command[1] == "exit") {
		t.processesInBackground++
		fmt.Fprintln(os.Stdout, fmt.Sprintf("[%d]\t", t.processesInBackground), strings.Join(command[1:], " "))
		fmt.Fprintln(os.Stdout, fmt.Sprintf("[%d]\t", t.processesInBackground), strings.Join(command[1:], " "), "\tDone")
		t.processesInBackground--
		return
	}

	if len(command) > 1 && command[1] == "echo" {
		t.processesInBackground++
		t.echo(command[1:])
		fmt.Fprintln(os.Stdout, fmt.Sprintf("[%d]\t", t.processesInBackground), strings.Join(command[1:], " "))
		fmt.Fprintln(os.Stdout, fmt.Sprintf("[%d]\t", t.processesInBackground), strings.Join(command[1:], " "), "\tDone")
		t.processesInBackground--
		return
	}

	goodToGo := make(chan bool)
	go t.executeInBackground(command[1:], goodToGo)
	<-goodToGo
}

func (t *terminal) executeInBackground(command []string, goodToGo chan<- bool) {
	// Catch SIGINT signals
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGINT)
	defer signal.Stop(sigint)

	if binary, err := exec.LookPath(command[0]); err == nil {
		attr := new(os.ProcAttr)
		attr.Dir, _ = os.Getwd()
		attr.Env = os.Environ()
		attr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}

		/*
		Start the process in its own process group by setting the Setpgid and Pgid fields in syscall.SysProcAttr,
		to prevent the process to receive signals sent directly to the parent process (gosh).
		By default, child processes start in the same process group, so without this piece of code,
		the process would be subject to signals, and would be killed if ctrl+c was pressed.
		*/
		attr.Sys = &syscall.SysProcAttr{
			Setpgid: true,
			Pgid:    0,
		}

		if process, err := os.StartProcess(binary, command, attr); err == nil {
			t.processesInBackground++
			processNumber := t.processesInBackground
			fmt.Fprintln(os.Stdout, fmt.Sprintf("[%d]\t%d\t", processNumber, process.Pid), strings.Join(command, " "))

			// Save in a map the command of this background process for an eventual foreground need
			t.backgroundProcesses[processNumber] = strings.Join(command, " ")

			// Release the main thread to continue executing
			goodToGo <- true

			// Channel to terminate the goroutine when the process finishes
			processFinished := make(chan bool)

			// Start a goroutine to handle signals coming in
			go func () {
				for {
					select {
					case <-sigint:
						// When a SIGINT signal arrives, only kill the process if foreground is active,
						// And if it is the last process to be ran.
						if t.fgActive && processNumber == t.processesInBackground {
							process.Kill()
						}

					// When the process finishes, terminate the goroutine
					case <-processFinished:
						return
					}
				}
			}()

			// Wait for the process to finish
			process.Wait()

			processFinished <- true

			fmt.Fprintln(os.Stdout, fmt.Sprintf("[%d]\t%d\t", processNumber, process.Pid), strings.Join(command, " "), "\tDone")
			t.processesInBackground--

			// Delete the command of this background process from the map
			delete(t.backgroundProcesses, processNumber)

			// Do not refresh the prompt if foreground is active
			if !t.fgActive {
				t.refresh()
			} else {
				// If foreground is active, release the waiting main thread
				t.waitBackgroundProcess <- true
			}

		} else {
			fmt.Fprintln(os.Stderr, err)
			// Release the main thread to continue executing
			goodToGo <- true
		}

	} else {
		fmt.Fprintln(os.Stderr, err)
		// Release the main thread to continue executing
		goodToGo <- true
	}
}

func (t *terminal) foreground() {
	if t.processesInBackground > 0 {
		// Print the command of the last background process ran
		fmt.Fprintln(os.Stdout, t.backgroundProcesses[t.processesInBackground])
		t.fgActive = true

		// Wait for the process to finishes
		<-t.waitBackgroundProcess

		// Clean the prompt
		t.line = t.line[:0]
		t.position = 0

	} else {
		fmt.Fprintln(os.Stderr, "No process running in background")
	}

	t.fgActive = false
}

func (t *terminal) cdCommand(command []string) error {
	if len(command) < 2 {
		return t.cd(os.Getenv("HOME"))
	} else {
		return t.cd(command[1])
	}
}

// Change directory.
func (t *terminal) cd(path string) error {
	currentDir, _ := os.Getwd()

	// If the informed path is a dash ("-") return to the last directory.
	if path == "-" && t.lastDirectory != "" {
		path = t.lastDirectory
		fmt.Fprintln(os.Stdout, path)
	}

	if err := os.Chdir(path); err != nil {
		fmt.Fprintln(os.Stderr, err)
		t.lastExitCode = 1
		return err
	}

	t.lastExitCode = 0
	t.lastDirectory = currentDir

	return nil
}

func (t *terminal) echo(command []string) error {
	if len(command) > 1 {
		commandStr := strings.Join(command[1:], " ")

		// Check if there is any "$" in the command
		// If it is, replace the $word for the correspondent env variable
		if strings.Contains(commandStr, "$") {
			// Do the checking word by word of the command
			for _, word := range command[1:] {
				words := strings.Split(word, "$")

				// Start evaluating from the second array element
				// Because, it only matters the string after the "$" symbol
				for _, envVariable := range words[1:] {
					var envVariableValue string

					// If the env var is $?
					// Get the exit code of the last executed process.
					if envVariable == "?" {
						envVariableValue = fmt.Sprint(t.lastExitCode)
					} else {
						envVariableValue = os.Getenv(envVariable)
					}

					commandStr = strings.ReplaceAll(commandStr, "$"+envVariable, envVariableValue)
				}
			}

			copy(command[1:], strings.Split(commandStr, " "))
		}
	}

	return t.execute(command)
}

func exit() {
	fmt.Fprintln(os.Stdout)
	os.Exit(0)
}
