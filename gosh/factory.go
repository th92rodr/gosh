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
	supported    bool
	originalMode syscall.Termios

	reader       *bufio.Reader
	nextInput	<-chan input
	winch        chan os.Signal
	sigint		 chan os.Signal
	columns      int

	history      stack
	currentLine

	eof bool

	process
}

type currentLine struct {
	line     []rune
	position int

	pendingEsc []rune
	escIsOn bool

	needRefresh bool

	historyPosition int
	historyEnd      string
}

type process struct {
	lastDirectory string
	lastExitCode int

	processesInBackground int
}

type input struct {
	char   rune
	err error
}

func New() *terminal {
	var terminal terminal
	terminal.reader = bufio.NewReader(os.Stdin)
	terminal.supported = isTerminalSupported()
	terminal.eof = false

	terminal.newStack()

	if terminal.supported {
		syscall.Syscall(syscall.SYS_IOCTL, uintptr(syscall.Stdin), syscall.TCGETS, uintptr(unsafe.Pointer(&terminal.originalMode)))
		mode := terminal.originalMode
		mode.Iflag &^= syscall.ICRNL | syscall.INPCK | syscall.ISTRIP | syscall.IXON
		mode.Cflag |= syscall.CS8
		mode.Lflag &^= syscall.ECHO | syscall.ICANON | syscall.IEXTEN
		syscall.Syscall(syscall.SYS_IOCTL, uintptr(syscall.Stdin), syscall.TCSETS, uintptr(unsafe.Pointer(&mode)))

		terminal.createWinchChannel()

		terminal.createSigIntChannel()
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

func (t *terminal) createSigIntChannel() {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGINT)
	t.sigint = sigint
}

func (t *terminal) close() {
	signal.Stop(t.winch)
	signal.Stop(t.sigint)

	if t.supported {
		syscall.Syscall(syscall.SYS_IOCTL, uintptr(syscall.Stdin), syscall.TCSETS, uintptr(unsafe.Pointer(&t.originalMode)))
	}
}
