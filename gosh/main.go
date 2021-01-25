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

	terminal.run("clear")

	for {
		input := terminal.prompt()

		if terminal.eof {
			break
		}

		if err := terminal.run(input); err != nil {
			break
		}
	}

	fmt.Fprintln(os.Stdout)
}
