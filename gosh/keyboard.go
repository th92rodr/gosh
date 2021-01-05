package main

import (
	"fmt"
	"os"
	"unicode"
)

/* Moving cursor */

// Move cursor to beginning of line
func (t *terminal) home() {
	t.position = 0
	t.needRefresh = true
}

// Move cursor to end of line
func (t *terminal) end() {
	t.position = len(t.line)
	t.needRefresh = true
}

// Move cursor one character left
func (t *terminal) left() {
	if t.position > 0 {
		t.position--
		t.needRefresh = true
	} else {
		doBeep()
	}
}

// Move cursor one character right
func (t *terminal) right() {
	if t.position < len(t.line) {
		t.position++
		t.needRefresh = true
	} else {
		doBeep()
	}
}

// Move cursor one word left
func (t *terminal) wordLeft() {
	if t.position > 0 {
		var spaceHere, spaceLeft, leftKnown bool
		for {
			t.position--
			if t.position == 0 {
				break
			}

			if leftKnown {
				spaceHere = spaceLeft
			} else {
				spaceHere = unicode.IsSpace(t.line[t.position])
			}

			spaceLeft, leftKnown = unicode.IsSpace(t.line[t.position-1]), true
			if !spaceHere && spaceLeft {
				break
			}
		}
	} else {
		doBeep()
	}
}

// Move cursor one word right
func (t *terminal) wordRight() {
	if t.position < len(t.line) {
		var spaceHere, spaceLeft, hereKnown bool
		for {
			t.position++
			if t.position == len(t.line) {
				break
			}

			if hereKnown {
				spaceLeft = spaceHere
			} else {
				spaceLeft = unicode.IsSpace(t.line[t.position-1])
			}

			spaceHere, hereKnown = unicode.IsSpace(t.line[t.position]), true
			if spaceHere && !spaceLeft {
				break
			}
		}
	} else {
		doBeep()
	}
}

/* Deleting character */

// Delete current character
func (t *terminal) delete() {
	if t.position >= len(t.line) {
		doBeep()
	} else {
		t.line = append(t.line[:t.position], t.line[t.position+1:]...)
	}
}

func (t *terminal) deleteNextWord() {
	if t.position == len(t.line) {
		doBeep()
		return
	}

	// Remove whitespace to the right
	var buf []rune // Store the deleted chars in a buffer
	for {
		if t.position == len(t.line) || !unicode.IsSpace(t.line[t.position]) {
			break
		}
		buf = append(buf, t.line[t.position])
		t.line = append(t.line[:t.position], t.line[t.position+1:]...)
	}

	// Remove non-whitespace to the right
	for {
		if t.position == len(t.line) || unicode.IsSpace(t.line[t.position]) {
			break
		}
		buf = append(buf, t.line[t.position])
		t.line = append(t.line[:t.position], t.line[t.position+1:]...)
	}
}

// Delete current word
func (t *terminal) eraseWord() {
	if t.position == 0 {
		doBeep()
		return
	}

	// Remove whitespace to the left
	var buf []rune // Store the deleted chars in a buffer
	for {
		if t.position == 0 || !unicode.IsSpace(t.line[t.position-1]) {
			break
		}
		buf = append(buf, t.line[t.position-1])
		t.line = append(t.line[:t.position-1], t.line[t.position:]...)
		t.position--
	}

	// Remove non-whitespace to the left
	for {
		if t.position == 0 || unicode.IsSpace(t.line[t.position-1]) {
			break
		}
		buf = append(buf, t.line[t.position-1])
		t.line = append(t.line[:t.position-1], t.line[t.position:]...)
		t.position--
	}

	// Invert the buffer and save the result on the killRing
	var newBuf []rune
	for i := len(buf) - 1; i >= 0; i-- {
		newBuf = append(newBuf, buf[i])
	}

	t.needRefresh = true
}

// Delete from start of line to cursor
func (t *terminal) ctrlU() {
	t.line = t.line[t.position:]
	t.position = 0
	t.needRefresh = true
}

// Delete from cursor to end of line
func (t *terminal) ctrlK() {
	if t.position >= len(t.line) {
		doBeep()
	} else {
		t.line = t.line[:t.position]
		t.needRefresh = true
	}
}

// Delete character before cursor
func (t *terminal) ctrlH() {
	if t.position <= 0 {
		doBeep()
	} else {
		t.line = append(t.line[:t.position-1], t.line[t.position:]...)
		t.position--
		t.needRefresh = true
	}
}

/* History */

func (t *terminal) up() {
	if t.historyPosition > 0 {
		if t.historyPosition == len(t.history) {
			t.historyEnd = string(t.line)
		}
		t.historyPosition--
		t.line = []rune(t.history[t.historyPosition])
		t.position = len(t.line)
		t.needRefresh = true
	} else {
		doBeep()
	}
}

func (t *terminal) down() {
	if t.historyPosition < len(t.history) {
		t.historyPosition++
		if t.historyPosition == len(t.history) {
			t.line = []rune(t.historyEnd)
		} else {
			t.line = []rune(t.history[t.historyPosition])
		}
		t.position = len(t.line)
		t.needRefresh = true
	} else {
		doBeep()
	}
}

/* */

// Transpose previous character with current character
func (t *terminal) ctrlT() {
	if len(t.line) < 2 || t.position < 1 {
		doBeep()
	} else {
		if t.position == len(t.line) {
			t.position--
		}
		t.line[t.position-1], t.line[t.position] = t.line[t.position], t.line[t.position-1]
		t.position++
		t.needRefresh = true
	}
}

// Clear screen
func (t *terminal) ctrlL() {
	t.eraseScreen()
	t.needRefresh = true
}

func (t *terminal) ctrlD() {
	if t.position == 0 && len(t.line) == 0 {
		t.eof = true
		return
	}

	if t.position >= len(t.line) {
		doBeep()
	} else {
		t.line = append(t.line[:t.position], t.line[t.position+1:]...)
		t.needRefresh = true
	}
}

func (t *terminal) ctrlC() {
	fmt.Fprintln(os.Stdout, "^C")
	t.line = t.line[:0]
	t.position = 0
	fmt.Fprint(os.Stdout, promptText)
}
