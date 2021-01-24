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

	run("clear")

	for {
		input := terminal.prompt()

		if terminal.eof {
			break
		}

		if err := run(input); err != nil {
			break
		}
	}

	fmt.Fprintln(os.Stdout)
}
