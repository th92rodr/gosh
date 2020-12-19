package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func (t *terminal) Prompt() (string, error) {
	fmt.Fprint(os.Stdout, promptText)

	t.getColumns()

	t.line = make([]rune, 0)
	t.position = 0
	t.historyPosition = len(t.history)
	t.historyEnd = ""

	t.startPrompt()

mainLoop:
	for {
		next, err := t.readNext()
		if err != nil {
			return "", err
		}

		switch value := next.(type) {
		case rune:
			switch value {
			case cr, lf:
				fmt.Fprintln(os.Stdout)
				break mainLoop
			case ctrlA: // Start of line
				t.ctrlA()
			case ctrlE: // End of line
				t.ctrlE()
			case ctrlB: // left
				t.ctrlB()
			case ctrlF: // right
				t.ctrlF()
			case ctrlD: // del
				if err := t.ctrlD(); err != nil {
					return "", err
				}
			case ctrlK: // delete remainder of line
				t.ctrlK()
			case ctrlP: // up
				t.ctrlP()
			case ctrlN: // down
				t.ctrlN()
			case ctrlT: // transpose prev rune with rune under cursor
				t.ctrlT()
			case ctrlL: // clear screen
				t.ctrlL()
			case ctrlH, bs: // Backspace
				t.ctrlH()
			case ctrlU: // Erase entire line
				t.ctrlU()
			// Catch unhandled control codes (anything <= 31)
			case 0, 3, 7, 9, 15:
				fallthrough
			case 17, 18, 19, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31:
				fmt.Print(beep)
			default:
				if t.position == len(t.line) && len(promptText)+len(t.line) < t.columns {
					t.line = append(t.line, value)
					fmt.Printf("%c", value)
					t.position++
				} else {
					t.line = append(t.line[:t.position], append([]rune{value}, t.line[t.position:]...)...)
					t.position++
					t.refresh(string(t.line), t.position)
				}
			}

		case action:
			switch value {
			case del:
				t.delete()
			case left:
				t.left()
			case right:
				t.right()
			case up:
				t.up()
			case down:
				t.down()
			}

			t.refresh(string(t.line), t.position)
		}
	}

	return string(t.line), nil
}

func (t *terminal) startPrompt() {
	next := make(chan nexter)
	go func() {
		for {
			var n nexter
			n.r, _, n.err = t.reader.ReadRune()
			next <- n

			// Shut down nexter loop when an end condition has been reached
			if n.err != nil || n.r == '\n' || n.r == '\r' || n.r == ctrlC || n.r == ctrlD {
				close(next)
				return
			}
		}
	}()
	t.next = next
}

func (t *terminal) readNext() (interface{}, error) {
	if len(t.pending) > 0 {
		rv := t.pending[0]
		t.pending = t.pending[1:]
		return rv, nil
	}

	var r rune

	select {
	case thing := <-t.next:
		if thing.err != nil {
			return nil, thing.err
		}
		r = thing.r
	case <-t.winch:
		t.getColumns()
		return winch, nil
	}

	if r != esc {
		return r, nil
	}
	t.pending = append(t.pending, r)

	timeout := time.After(50 * time.Millisecond)
	flag, err := t.nextPending(timeout)
	if err != nil {
		if err == timedOut {
			return flag, nil
		}
		return unknown, err
	}

	switch flag {
	case '[':
		code, err := t.nextPending(timeout)
		if err != nil {
			if err == timedOut {
				return code, nil
			}
			return unknown, err
		}
		switch code {
		case 'A':
			t.pending = t.pending[:0] // escape code complete
			return up, nil
		case 'B':
			t.pending = t.pending[:0] // escape code complete
			return down, nil
		case 'C':
			t.pending = t.pending[:0] // escape code complete
			return right, nil
		case 'D':
			t.pending = t.pending[:0] // escape code complete
			return left, nil
		case 'Z':
			t.pending = t.pending[:0] // escape code complete
			return shiftTab, nil
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			num := []rune{code}
			for {
				code, err := t.nextPending(timeout)
				if err != nil {
					if err == timedOut {
						return code, nil
					}
					return nil, err
				}
				if code < '0' || code > '9' {
					if code != '~' {
						// escape code went off the rails
						rv := t.pending[0]
						t.pending = t.pending[1:]
						return rv, nil
					}
					break
				}
				num = append(num, code)
			}
			t.pending = t.pending[:0] // escape code complete
			x, _ := strconv.ParseInt(string(num), 10, 32)
			switch x {
			case 2:
				return insert, nil
			case 3:
				return del, nil
			case 5:
				return pageUp, nil
			case 6:
				return pageDown, nil
			case 15:
				return f5, nil
			case 17:
				return f6, nil
			case 18:
				return f7, nil
			case 19:
				return f8, nil
			case 20:
				return f9, nil
			case 21:
				return f10, nil
			case 23:
				return f11, nil
			case 24:
				return f12, nil
			default:
				return unknown, nil
			}
		}

	case 'O':
		code, err := t.nextPending(timeout)
		if err != nil {
			if err == timedOut {
				return code, nil
			}
			return nil, err
		}
		t.pending = t.pending[:0] // escape code complete
		switch code {
		case 'H':
			return home, nil
		case 'F':
			return end, nil
		case 'P':
			return f1, nil
		case 'Q':
			return f2, nil
		case 'R':
			return f3, nil
		case 'S':
			return f4, nil
		default:
			return unknown, nil
		}
	default:
		rv := t.pending[0]
		t.pending = t.pending[1:]
		return rv, nil
	}

	return nil, nil
}

func (t *terminal) nextPending(timeout <-chan time.Time) (rune, error) {
	select {
	case thing := <-t.next:
		if thing.err != nil {
			return 0, thing.err
		}
		t.pending = append(t.pending, thing.r)
		return thing.r, nil
	case <-timeout:
		rv := t.pending[0]
		t.pending = t.pending[1:]
		return rv, timedOut
	}
}
