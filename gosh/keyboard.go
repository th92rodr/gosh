package main

import (
	"fmt"
	"os"
)

//// Moving cursor

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
		for {
			t.position--

			// Check for begining of line
			if t.position == 0 {
				break
			}

			spaceHere := isABlankSpace(t.line[t.position])
			spaceLeft := isABlankSpace(t.line[t.position-1])

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
		for {
			t.position++

			// Check for end of line
			if t.position == len(t.line) {
				break
			}

			spaceHere := isABlankSpace(t.line[t.position])
			spaceLeft := isABlankSpace(t.line[t.position-1])

			if spaceHere && !spaceLeft {
				break
			}
		}
	} else {
		doBeep()
	}
}

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

//// Deleting characters

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
	for {
		if t.position == len(t.line) || !isABlankSpace(t.line[t.position]) {
			break
		}
		t.line = append(t.line[:t.position], t.line[t.position+1:]...)
	}

	// Remove non-whitespace to the right
	for {
		if t.position == len(t.line) || isABlankSpace(t.line[t.position]) {
			break
		}
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
	for {
		if t.position == 0 || !isABlankSpace(t.line[t.position-1]) {
			break
		}
		t.line = append(t.line[:t.position-1], t.line[t.position:]...)
		t.position--
	}

	// Remove non-whitespace to the left
	for {
		if t.position == 0 || isABlankSpace(t.line[t.position-1]) {
			break
		}
		t.line = append(t.line[:t.position-1], t.line[t.position:]...)
		t.position--
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

//// History

// Previous command from history
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

// Next command from history
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

////

// Clear screen
func (t *terminal) ctrlL() {
	t.eraseScreen()
	t.needRefresh = true
}

// End of File - if line is empty quits application
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

// Reset input
func (t *terminal) ctrlC() {
	fmt.Fprintln(os.Stdout, "^C")
	t.line = t.line[:0]
	t.position = 0
	fmt.Fprint(os.Stdout, promptText)
}
