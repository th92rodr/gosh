package main

import (
	"fmt"
	"os"
)

func main() {
	terminal := New()
	defer terminal.close()

	if !terminal.supported {
		return
	}

	run(parseInput("clear"))

	for {
		input := terminal.prompt()

		if terminal.eof {
			break
		}

		command := parseInput(input)
		if err := run(command); err != nil {
			break
		}
	}

	fmt.Fprintln(os.Stdout)
}
