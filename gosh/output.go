package main

import (
	"fmt"
	"os"
	"syscall"
	"unicode/utf8"
	"unsafe"
)

func (t *terminal) refresh(buf string, position int) error {
	t.cursorPosition(0)

	fmt.Fprint(os.Stdout, promptText)

	pLen := utf8.RuneCountInString(promptText)
	bLen := utf8.RuneCountInString(buf)

	if pLen+bLen <= t.columns {
		fmt.Fprint(os.Stdout, buf)
		t.eraseLine()
		t.cursorPosition(pLen + position)
	} else {
		// Find space available
		space := t.columns - pLen
		space-- // space for cursor
		start := position - space/2
		end := start + space

		if end > bLen {
			end = bLen
			start = end - space
		}

		if start < 0 {
			start = 0
			end = space
		}
		position -= start

		// Leave space for markers
		if start > 0 {
			start++
		}

		if end < bLen {
			end--
		}

		line := []rune(buf)
		line = line[start:end]

		// Output
		if start > 0 {
			fmt.Fprint(os.Stdout, "{")
		}

		fmt.Fprint(os.Stdout, string(line))
		if end < bLen {
			fmt.Fprint(os.Stdout, "}")
		}

		// Set cursor position
		t.eraseLine()
		t.cursorPosition(pLen + position)
	}

	return nil
}

type winSize struct {
	row, col       uint16
	xpixel, ypixel uint16
}

func (t *terminal) getColumns() {
	var ws winSize
	if ok, _, _ := syscall.Syscall(syscall.SYS_IOCTL, uintptr(syscall.Stdout), syscall.TIOCGWINSZ, uintptr(unsafe.Pointer(&ws))); ok < 0 {
		t.columns = 80
	}
	t.columns = int(ws.col)
}

func (t *terminal) cursorPosition(position int) {
	fmt.Printf("\x1b[%dG", position+1)
}

func (t *terminal) eraseLine() {
	fmt.Print("\x1b[0K")
}

func (t *terminal) eraseScreen() {
	fmt.Print("\x1b[H\x1b[2J")
}

func doBeep() {
	fmt.Fprint(os.Stdout, beep)
}
