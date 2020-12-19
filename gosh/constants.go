package main

import (
	"errors"
)

const promptText = "gosh >  "

type action int

const (
	left action = iota
	right
	up
	down
	home
	end
	insert
	del
	pageUp
	pageDown
	f1
	f2
	f3
	f4
	f5
	f6
	f7
	f8
	f9
	f10
	f11
	f12
	shiftTab
	winch
	unknown
)

const (
	ctrlA = 1
	ctrlB = 2
	ctrlC = 3
	ctrlD = 4
	ctrlE = 5
	ctrlF = 6
	ctrlH = 8
	tab   = 9
	lf    = 10 // Line Feed: Causes the cursor to jump to the next line
	ctrlK = 11
	ctrlL = 12
	cr    = 13 // Carriage Return: Moves the cursor back to the first position of the line
	ctrlN = 14
	ctrlP = 16
	ctrlT = 20
	ctrlU = 21
	esc   = 27
	bs    = 127 // Backspace: Lets the cursor move back one step

	beep = "\a"
)

var timedOut = errors.New("timeout")
