package main

import (
	"strings"
)

// Parse input and return a slice of commands,
// which one in the following format [command, args...]
func parseInput(input string) [][]string {
	// Remove the newline character.
	input = strings.TrimSuffix(input, "\n")

	// Split the input to separate the command and the arguments.
	args := strings.Split(input, " ")

	if len(args) > 1 {
		// Remove empty elements.
		args = removeEmptyElements(args)
	}

	var commands [][]string

	// Separate the commands by logic operators (AND and OR),
	// in the following format [command 1] [AND] [command 2] [OR] [command 3]
	startIndex := 0
	for index, arg := range args {
		if arg == andOperator || arg == orOperator || arg == semiColonOperator {
			commands = append(commands, args[startIndex:index])
			commands = append(commands, args[index:index+1])
			startIndex = index + 1
		}
	}
	commands = append(commands, args[startIndex:])

	return commands
}

func removeEmptyElements(slice []string) []string {
	var newSlice []string

	for _, str := range slice {
		if str != "" {
			newSlice = append(newSlice, str)
		}
	}

	return newSlice
}
