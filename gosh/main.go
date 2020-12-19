package main

import (
	"fmt"
)

func main() {
	terminal := New()
	defer terminal.Close()

	if !terminal.supported {
		return
	}

	for {
		if input, err := terminal.Prompt(); err == nil {
			command := ParseInput(input)
			if command[0] == "exit" {
				break
			}
			Run(command)
		} else {
			break
		}
	}

}
