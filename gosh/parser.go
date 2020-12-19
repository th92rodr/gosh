package main

import (
	"strings"
)

// parse input and return in the following format [command, args...]
func ParseInput(input string) []string {
	// Remove the newline character.
	input = strings.TrimSuffix(input, "\n")

	// Split the input to separate the command and the arguments.
	args := strings.Split(input, " ")

	if len(args) > 1 {
		// Remove empty elements.
		args = removeEmptyElements(args)
	}

	return args
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
