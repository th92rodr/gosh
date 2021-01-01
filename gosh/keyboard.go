package main

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

func (t *terminal) delete() {
	if t.position >= len(t.line) {
		doBeep()
	} else {
		t.line = append(t.line[:t.position], t.line[t.position+1:]...)
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
