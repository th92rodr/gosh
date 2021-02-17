package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"unicode/utf8"
)

type stack []string

func (t *terminal) initHistory() {
	t.history = make([]string, 0)

	if file, err := os.Open(t.historyFilename); err == nil {
		t.readHistory(file)
		file.Close()
	}
}

func (t *terminal) terminateHistory() {
	if file, err := os.Create(t.historyFilename); err == nil {
		t.writeHistory(file)
		file.Close()
	}
}

func (t *terminal) push(input string) {
	t.history = append(t.history, input)
}

func (t *terminal) top() string {
	return t.history[len(t.history)-1]
}

func (t *terminal) writeHistory(file io.Writer) error {
	for _, command := range t.history {
		if _, err := fmt.Fprintln(file, command); err != nil {
			return err
		}
	}

	return nil
}

func (t *terminal) readHistory(file io.Reader) error {
	lineNumber := 0
	reader := bufio.NewReader(file)

	for {
		line, tooLong, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if tooLong {
			return errors.New(fmt.Sprintf("line %d is too long", lineNumber+1))
		}
		if !utf8.Valid(line) {
			return errors.New(fmt.Sprintf("invalid string at line %d", lineNumber+1))
		}

		lineNumber++
		t.history = append(t.history, string(line))
	}

	return nil
}
