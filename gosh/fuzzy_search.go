package main

import (
	"unicode/utf8"
)

func find(source string, targets []string) []string {
	var matches []string

	if source == "" {
		return matches
	}

	for index := len(targets) - 1; index >= 0; index-- {
		if match(source, targets[index]) {
			matches = append(matches, targets[index])
		}
	}

	return matches
}

func match(source, target string) bool {
	lenDiff := len(target) - len(source)

	if lenDiff < 0 {
		return false
	}

	if lenDiff == 0 && source == target {
		return true
	}

MatchLoop:
	for _, s := range source {
		for index, t := range target {
			if s == t {
				target = target[index+utf8.RuneLen(t):]
				continue MatchLoop
			}
		}
		return false
	}

	return true
}
