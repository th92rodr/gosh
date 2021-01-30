package main

import (
    "bufio"
    "fmt"
    "os"
    "os/exec"
    "strings"
)

func main() {
    reader := bufio.NewReader(os.Stdin)

    for {
        fmt.Print("sh >  ")

        // Read the keyboard input.
        input, err := reader.ReadString('\n')
        if err != nil {
            fmt.Fprintln(os.Stderr, err)
        }

        // Handle the execution of the input.
        if err = execInput(input); err != nil {
            fmt.Fprintln(os.Stderr, err)
        }
    }
}

func execInput(input string) error {
    // Remove the newline character.
    input = strings.TrimSuffix(input, "\n")

    // Split the input to separate the command and the arguments.
    args := strings.Split(input, " ")

    // Check for built-in commands.
    switch args[0] {
    case "":
        return nil
    case "cd":
        // 'cd' to home dir when path is empty
        if len(args) < 2 {
            homeDir, _ := os.UserHomeDir()
            return os.Chdir(homeDir)
        }
        // Change the directory and return the error.
        return os.Chdir(args[1])
    case "exit":
        os.Exit(0)
    }

    // Prepare the command to execute.
    cmd := exec.Command(args[0], args[1:]...)

    // Set the correct output device.
    cmd.Stderr = os.Stderr
    cmd.Stdout = os.Stdout

    // Execute the command and return the error.
    return cmd.Run()
}
