package main

import (
	"fmt"
	"os"
	"time"
)

func (t *terminal) prompt() string {
	fmt.Fprint(os.Stdout, promptText)

	t.getColumns()

	t.line = make([]rune, 0)
	t.position = 0

	t.pendingEsc = make([]rune, 0)
	t.escIsOn = false

	t.historyPosition = len(t.history)
	t.historyEnd = ""

	t.ctrlRSearches = 0

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

			} else if isKeyStroke(char) {
				t.pendingEsc = append(t.pendingEsc, char)
				break mainLoop

			} else {
				if t.position == len(t.line) && len(promptText)+len(t.line) < t.columns {
					t.line = append(t.line, char)
					fmt.Fprintf(os.Stdout, "%s", string(char))
					t.position++
				} else {
					t.line = append(t.line[:t.position], append([]rune{char}, t.line[t.position:]...)...)
					t.position++
					t.refresh()
				}
			}

		case <-timeout:
			break mainLoop

		case <-t.winch:
			t.getColumns()
			break mainLoop

		case <-t.sigint:
			t.pendingEsc = append(t.pendingEsc, CTRL_C[0])
			break mainLoop
		}
	}

	if len(t.pendingEsc) > 0 {
		t.executeEscapeKey()
		if !t.eof {
			goto mainLoop
		}
	}

	if len(t.line) > 0 {
		t.push(string(t.line))
		return string(t.line)
	}

	return ""
}

func (t *terminal) startPrompt() {
	next := make(chan input)

	// Keep reading inputs until an end condition is reached.
	go func() {
		for {
			var i input
			i.char, _, i.err = t.reader.ReadRune()
			next <- i

			// Stop next input loop when an end condition has been reached.
			if i.err != nil || i.char == '\n' || i.char == '\r' || i.char == ENTER {
				close(next)
				return
			}
		}
	}()

	t.nextInput = next
}

func (t *terminal) executeEscapeKey() {
	for index, key := range keys {
		if ok := isEqual(t.pendingEsc, key); ok {
			switch keysArrayIndexMapsToKeyName[index] {
			case "home", "ctrlA":
				t.home()
			case "end", "ctrlE":
				t.end()
			case "left", "ctrlB":
				t.left()
			case "right", "ctrlF":
				t.right()
			case "altB", "wordLeft":
				t.wordLeft()
			case "altF", "wordRight":
				t.wordRight()
			case "ctrlT":
				t.ctrlT()

			case "delete":
				t.delete()
			case "altD":
				t.deleteNextWord()
			case "altBackspace", "ctrlW":
				t.eraseWord()
			case "ctrlU":
				t.ctrlU()
			case "ctrlK":
				t.ctrlK()
			case "ctrlH", "backspace":
				t.ctrlH()

			case "up", "ctrlP":
				t.up()
			case "down", "ctrlN":
				t.down()
			case "ctrlR":
				t.ctrlR()

			case "ctrlL":
				t.ctrlL()
			case "ctrlD":
				t.ctrlD()
			case "ctrlC":
				t.ctrlC()

			case "tab":
				t.tabCompleter()
			}

			if t.needRefresh {
				t.refresh()
			}

			t.pendingEsc = t.pendingEsc[:0]		// Clean slice
			t.escIsOn = false
			return
		}
	}

	t.pendingEsc = t.pendingEsc[:0]		// Clean slice
	t.escIsOn = false
}

func isKeyStroke(char rune) bool {
	return char == CTRL_A[0] ||
		char == CTRL_B[0] ||
		char == CTRL_D[0] ||
		char == CTRL_E[0] ||
		char == CTRL_F[0] ||
		char == CTRL_G[0] ||
		char == CTRL_H[0] ||
		char == CTRL_K[0] ||
		char == CTRL_L[0] ||
		char == CTRL_N[0] ||
		char == CTRL_O[0] ||
		char == CTRL_P[0] ||
		char == CTRL_Q[0] ||
		char == CTRL_R[0] ||
		char == CTRL_S[0] ||
		char == CTRL_T[0] ||
		char == CTRL_U[0] ||
		char == CTRL_V[0] ||
		char == CTRL_W[0] ||
		char == CTRL_X[0] ||
		char == CTRL_Y[0] ||
		char == TAB[0] ||
		char == LINE_FEED[0] ||
		char == CARRIAGE_RETURN[0] ||
		char == BACKSPACE[0]
}
