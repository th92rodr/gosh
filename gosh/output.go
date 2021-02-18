package main

import (
	"fmt"
	"os"
	"syscall"
	"unicode/utf8"
	"unsafe"
)

func (t *terminal) refresh() {
	if t.isMultilineCommand {
		t.cursorPosition(0)
		fmt.Fprint(os.Stdout, string(t.line))
		t.eraseLine()
		t.cursorPosition(t.position)
		return
	}

	t.cursorPosition(0)

	fmt.Fprint(os.Stdout, promptText)

	pLen := utf8.RuneCountInString(promptText)
	bLen := utf8.RuneCountInString(string(t.line))

	if pLen+bLen <= t.columns {
		fmt.Fprint(os.Stdout, string(t.line))
		t.eraseLine()
		t.cursorPosition(pLen + t.position)

	} else {
		// Find space available
		space := t.columns - pLen
		space-- // space for cursor
		start := t.position - space/2
		end := start + space

		if end > bLen {
			end = bLen
			start = end - space
		}

		if start < 0 {
			start = 0
			end = space
		}
		t.position -= start

		// Leave space for markers
		if start > 0 {
			start++
		}

		if end < bLen {
			end--
		}

		line := t.line[start:end]

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
		t.cursorPosition(pLen + t.position)
	}
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
