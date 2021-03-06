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
	historyFilename string

	currentLine

	eof bool

	process
}

type currentLine struct {
	line     []rune
	position int

	isMultilineCommand bool
	multiline [][]rune // Each line of the multiline command will be in a separated slice

	pendingEsc []rune
	escIsOn bool

	needRefresh bool

	historyPosition int
	historyEnd      string

	ctrlRSearches int
}

type process struct {
	lastDirectory string
	lastExitCode int

	processesInBackground int	// Number of process running in background
	backgroundProcesses map[int]string // Commands of the background processes ran
	fgActive bool
	waitBackgroundProcess chan bool
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

	terminal.historyFilename = ".history"
	terminal.initHistory()

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

	terminal.lastDirectory = ""
	terminal.lastExitCode = 0

	terminal.processesInBackground = 0
	terminal.backgroundProcesses = make(map[int]string)
	terminal.fgActive = false
	terminal.waitBackgroundProcess = make(chan bool)

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

	t.terminateHistory()
}
