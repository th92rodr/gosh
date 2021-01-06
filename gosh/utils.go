package main

import (
    "unicode"
)

func isEmpty(slice []rune) bool {
	return len(slice) == 0
}

func isABlankSpace(char rune) bool {
	return unicode.IsSpace(char)
}

// Compare two rune slices
func isEqual(a, b []rune) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}

	return true
}
