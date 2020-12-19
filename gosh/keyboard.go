package main

import (
	"fmt"
	"io"
)

func (t *terminal) ctrlA() {
	t.position = 0
	t.refresh(string(t.line), t.position)
}

func (t *terminal) ctrlE() {
	t.position = len(t.line)
	t.refresh(string(t.line), t.position)
}

func (t *terminal) ctrlB() {
	if t.position > 0 {
		t.position--
		t.refresh(string(t.line), t.position)
	} else {
		fmt.Print(beep)
	}
}

func (t *terminal) ctrlF() {
	if t.position < len(t.line) {
		t.position++
		t.refresh(string(t.line), t.position)
	} else {
		fmt.Print(beep)
	}
}

func (t *terminal) ctrlD() error {
	if t.position == 0 && len(t.line) == 0 {
		return io.EOF
	}

	if t.position >= len(t.line) {
		fmt.Print(beep)
	} else {
		t.line = append(t.line[:t.position], t.line[t.position+1:]...)
		t.refresh(string(t.line), t.position)
	}

	return nil
}

func (t *terminal) ctrlK() {
	if t.position >= len(t.line) {
		fmt.Print(beep)
	} else {
		t.line = t.line[:t.position]
		t.refresh(string(t.line), t.position)
	}
}

func (t *terminal) ctrlP() {
	if t.historyPosition > 0 {
		if t.historyPosition == len(t.history) {
			t.historyEnd = string(t.line)
		}
		t.historyPosition--
		t.line = []rune(t.history[t.historyPosition])
		t.position = len(t.line)
		t.refresh(string(t.line), t.position)
	} else {
		fmt.Print(beep)
	}
}

func (t *terminal) ctrlN() {
	if t.historyPosition < len(t.history) {
		t.historyPosition++
		if t.historyPosition == len(t.history) {
			t.line = []rune(t.historyEnd)
		} else {
			t.line = []rune(t.history[t.historyPosition])
		}
		t.position = len(t.line)
		t.refresh(string(t.line), t.position)
	} else {
		fmt.Print(beep)
	}
}

func (t *terminal) ctrlT() {
	if len(t.line) < 2 || t.position < 1 {
		fmt.Print(beep)
	} else {
		if t.position == len(t.line) {
			t.position--
		}
		t.line[t.position-1], t.line[t.position] = t.line[t.position], t.line[t.position-1]
		t.position++
		t.refresh(string(t.line), t.position)
	}
}

func (t *terminal) ctrlL() {
	t.eraseScreen()
	t.refresh(string(t.line), t.position)
}

func (t *terminal) ctrlH() {
	if t.position <= 0 {
		fmt.Print(beep)
	} else {
		t.line = append(t.line[:t.position-1], t.line[t.position:]...)
		t.position--
		t.refresh(string(t.line), t.position)
	}
}

func (t *terminal) ctrlU() {
	t.line = t.line[:0]
	t.position = 0
	t.refresh(string(t.line), t.position)
}

/////////////////

func (t *terminal) delete() {
	if t.position >= len(t.line) {
		fmt.Print(beep)
	} else {
		t.line = append(t.line[:t.position], t.line[t.position+1:]...)
	}
}

func (t *terminal) left() {
	if t.position > 0 {
		t.position--
	} else {
		fmt.Print(beep)
	}
}

func (t *terminal) right() {
	if t.position < len(t.line) {
		t.position++
	} else {
		fmt.Print(beep)
	}
}

func (t *terminal) up() {
	if t.historyPosition > 0 {
		if t.historyPosition == len(t.history) {
			t.historyEnd = string(t.line)
		}
		t.historyPosition--
		t.line = []rune(t.history[t.historyPosition])
		t.position = len(t.line)
	} else {
		fmt.Print(beep)
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
	} else {
		fmt.Print(beep)
	}
}
