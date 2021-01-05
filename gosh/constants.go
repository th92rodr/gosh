package main

const (
	promptText = "gosh >  "
	beep = "\a"
)

var keys = [][]rune{}

var keysArrayIndexMapsToKeyName = map[int]string {}

var keyNameMapsToEscSequence = map[string][]rune {
	"home": HOME,
	"end": END,

	"pageUp": PAGE_UP,
	"pageDown": PAGE_DOWN,

	"up": UP,
	"down": DOWN,
	"right": RIGHT,
	"left": LEFT,

	"insert": INSERT,
	"delete": DELETE,

	"f1": F1,
	"f2": F2,
	"f3": F3,
	"f4": F4,
	"f5": F5,
	"f6": F6,
	"f7": F7,
	"f8": F8,
	"f9": F9,
	"f10": F10,
	"f11": F11,
	"f12": F12,

	"altB": ALT_B,
	"altF": ALT_F,
	"altD": ALT_D,
	"altBackspace": ALT_BACKSPACE,

	"wordLeft": WORD_LEFT,
	"wordRight": WORD_RIGHT,

	"ctrlA": CTRL_A,
	"ctrlB": CTRL_B,
	"ctrlC": CTRL_C,
	"ctrlD": CTRL_D,
	"ctrlE": CTRL_E,
	"ctrlF": CTRL_F,
	"ctrlH": CTRL_H,
	"ctrlK": CTRL_K,
	"ctrlL": CTRL_L,
	"ctrlN": CTRL_N,
	"ctrlP": CTRL_P,
	"ctrlT": CTRL_T,
	"ctrlU": CTRL_U,
	"ctrlW": CTRL_W,

	"tab": TAB,
	"lineFeed": LINE_FEED,
	"carriageReturn": CARRIAGE_RETURN,
	"backspace": BACKSPACE,
}

var HOME = []rune{27, 91, 72}
var END = []rune{27, 91, 70}

var PAGE_UP = []rune{27, 91, 53, 126}
var PAGE_DOWN = []rune{27, 91, 54, 126}

var UP = []rune{27, 91, 65}
var DOWN = []rune{27, 91, 66}
var RIGHT = []rune{27, 91, 67}
var LEFT = []rune{27, 91, 68}

var INSERT = []rune{27, 91, 50, 126}
var DELETE = []rune{27, 91, 51, 126}

var F1 = []rune{27, 79, 80}
var F2 = []rune{27, 79, 81}
var F3 = []rune{27, 79, 82}
var F4 = []rune{27, 79, 83}
var F5 = []rune{27, 91, 49, 53, 126}
var F6 = []rune{27, 91, 49, 55, 126}
var F7 = []rune{27, 91, 49, 56, 126}
var F8 = []rune{27, 91, 49, 57, 126}
var F9 = []rune{27, 91, 50, 48, 126}
var F10 = []rune{}
var F11 = []rune{}
var F12 = []rune{27, 91, 50, 52, 126}

var ALT_B = []rune{27, 98}
var ALT_F = []rune{27, 102}
var ALT_D = []rune{27, 100}
var ALT_BACKSPACE = []rune{27, 127}

var WORD_LEFT = []rune{27, 91, 49, 59, 53, 68}		// CTRL + ARROW LEFT
var WORD_RIGHT = []rune{27, 91, 49, 59, 53, 67}		// CTRL + ARROW RIGHT

var CTRL_A = []rune{1}
var CTRL_B = []rune{2}
var CTRL_C = []rune{3}
var CTRL_D = []rune{4}
var CTRL_E = []rune{5}
var CTRL_F = []rune{6}
var CTRL_H = []rune{8}
var CTRL_K = []rune{11}
var CTRL_L = []rune{12}
var CTRL_N = []rune{14}
var CTRL_P = []rune{16}
var CTRL_T = []rune{20}
var CTRL_U = []rune{21}
var CTRL_W = []rune{23}

var TAB = []rune{9}
var LINE_FEED = []rune{10}
var CARRIAGE_RETURN = []rune{13}

var BACKSPACE = []rune{127}

const ESC = 27
const ENTER = 13

func init() {
	keys = make([][]rune, 0)
	index := 0

	for key, value := range keyNameMapsToEscSequence {
		keys = append(keys, value)
		keysArrayIndexMapsToKeyName[index] = key
		index++
	}
}
