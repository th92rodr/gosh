package main

import (
	"fmt"
	"os"
	"time"
)

func (t *terminal) prompt() (string, error) {
	fmt.Fprint(os.Stdout, promptText)

	t.getColumns()

	t.line = make([]rune, 0)
	t.position = 0

	t.pendingEsc = make([]rune, 0)
	t.escIsOn = false

	t.historyPosition = len(t.history)
	t.historyEnd = ""

	t.startPrompt()

	var timeout <-chan time.Time

mainLoop:
	for {
		select {
		case input, ok := <-t.nextInput:
			if !ok || input.err != nil || input.char == ENTER {
				break mainLoop
			}

			char := input.char

			if char == ESC || t.escIsOn {
				t.escIsOn = true
				t.pendingEsc = append(t.pendingEsc, char)

				// Wait for the rest of the escape sequence during 50 ms
				// If nothing else arrives, it was an actual ESC key
				timeout = time.After(50 * time.Millisecond)
			} else {
				if t.position == len(t.line) && len(promptText)+len(t.line) < t.columns {
					t.line = append(t.line, char)
					fmt.Printf("%s", string(char))
					t.position++
				} else {
					t.line = append(t.line[:t.position], append([]rune{char}, t.line[t.position:]...)...)
					t.position++
					t.refresh(string(t.line), t.position)
				}
			}

		case <-timeout:
			break mainLoop
		}
	}

	if len(t.pendingEsc) > 0 {
		if t.executeEscapeKey() {
			goto mainLoop
		}
	}

	if len(t.line) > 0 {
		return string(t.line), nil
	}

	return "", nil
}

func (t *terminal) startPrompt() {
	next := make(chan input)

	// keep reading inputs until an end condition is reached
	go func() {
		for {
			var i input
			i.char, _, i.err = t.reader.ReadRune()
			next <- i

			// Stop next input loop when an end condition has been reached
			if i.err != nil || i.char == '\n' || i.char == '\r' || i.char == ENTER {
				close(next)
				return
			}
		}
	}()

	t.nextInput = next
}

func (t *terminal) executeEscapeKey() bool {
	for index, key := range keys {
		if ok := isEqual(t.pendingEsc, key); ok {
			switch keysArrayIndexMapsToKeyName[index] {
			case "home", "ctrlA":
				t.home()
			case "end", "ctrlE":
				t.end()
			case "right", "ctrlF":
				t.right()
			case "left", "ctrlB":
				t.left()
			case "delete":
				t.delete()
			}

			if t.needRefresh {
				t.refresh(string(t.line), t.position)
			}

			t.pendingEsc = t.pendingEsc[:0]
			t.escIsOn = false

			return true
		}
	}

	return false
}

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
