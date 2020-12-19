package main

import (
	"bufio"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"unsafe"
)

type terminal struct {
	reader       *bufio.Reader
	supported    bool
	originalMode syscall.Termios
	next         <-chan nexter
	winch        chan os.Signal
	pending      []rune
	history      []string
	columns      int
	currentLine
}

type currentLine struct {
	line     []rune
	position int
	columns  int

	historyPosition int
	historyEnd      string
}

type nexter struct {
	r   rune
	err error
}

func New() *terminal {
	var terminal terminal
	terminal.reader = bufio.NewReader(os.Stdin)
	terminal.supported = isTerminalSupported()

	if terminal.supported {
		syscall.Syscall(syscall.SYS_IOCTL, uintptr(syscall.Stdin), syscall.TCGETS, uintptr(unsafe.Pointer(&terminal.originalMode)))
		mode := terminal.originalMode
		mode.Iflag &^= syscall.ICRNL | syscall.INPCK | syscall.ISTRIP | syscall.IXON
		mode.Cflag |= syscall.CS8
		mode.Lflag &^= syscall.ECHO | syscall.ICANON | syscall.IEXTEN
		syscall.Syscall(syscall.SYS_IOCTL, uintptr(syscall.Stdin), syscall.TCSETS, uintptr(unsafe.Pointer(&mode)))

		terminal.createWinchChannel()
	}

	return &terminal
}

func isTerminalSupported() bool {
	notSupported := map[string]bool{"": true, "dumb": true, "cons25": true}
	return !notSupported[strings.ToLower(os.Getenv("TERM"))]
}

func (t *terminal) createWinchChannel() {
	winch := make(chan os.Signal, 1)
	signal.Notify(winch, syscall.SIGWINCH)
	t.winch = winch
}

func (t *terminal) Close() {
	if t.supported {
		syscall.Syscall(syscall.SYS_IOCTL, uintptr(syscall.Stdin), syscall.TCSETS, uintptr(unsafe.Pointer(&t.originalMode)))
	}
}
