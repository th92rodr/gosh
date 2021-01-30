package main

const (
	promptText = "gosh >  "
	beep = "\a"

	andOperator = "&&"
	orOperator = "||"
	semiColonOperator = ";"
	backgroundOperator = "&"
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
	"altD": ALT_D,
	"altF": ALT_F,
	"altY": ALT_Y,
	"altBackspace": ALT_BACKSPACE,

	"wordLeft": WORD_LEFT,
	"wordRight": WORD_RIGHT,

	"ctrlA": CTRL_A,
	"ctrlB": CTRL_B,
	"ctrlC": CTRL_C,
	"ctrlD": CTRL_D,
	"ctrlE": CTRL_E,
	"ctrlF": CTRL_F,
	"ctrlG": CTRL_G,
	"ctrlH": CTRL_H,
	"ctrlK": CTRL_K,
	"ctrlL": CTRL_L,
	"ctrlN": CTRL_N,
	"ctrlO": CTRL_O,
	"ctrlP": CTRL_P,
	"ctrlQ": CTRL_Q,
	"ctrlR": CTRL_R,
	"ctrlS": CTRL_S,
	"ctrlT": CTRL_T,
	"ctrlU": CTRL_U,
	"ctrlV": CTRL_V,
	"ctrlW": CTRL_W,
	"ctrlX": CTRL_X,
	"ctrlY": CTRL_Y,
	"ctrlZ": CTRL_Z,

	"tab": TAB,
	"lineFeed": LINE_FEED,
	"carriageReturn": CARRIAGE_RETURN,
	"backspace": BACKSPACE,

	"shiftTab": SHIFT_TAB,
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
var F10 = []rune{27, 91, 50, 50, 126}
var F11 = []rune{27, 91, 50, 51, 126}
var F12 = []rune{27, 91, 50, 52, 126}

var ALT_B = []rune{27, 98}
var ALT_D = []rune{27, 100}
var ALT_F = []rune{27, 102}
var ALT_Y = []rune{27, 121}
var ALT_BACKSPACE = []rune{27, 127}

var WORD_LEFT = []rune{27, 91, 49, 59, 53, 68}		// CTRL + ARROW LEFT
var WORD_RIGHT = []rune{27, 91, 49, 59, 53, 67}		// CTRL + ARROW RIGHT

var CTRL_A = []rune{1}
var CTRL_B = []rune{2}
var CTRL_C = []rune{3}
var CTRL_D = []rune{4}
var CTRL_E = []rune{5}
var CTRL_F = []rune{6}
var CTRL_G = []rune{7}
var CTRL_H = []rune{8}
var CTRL_K = []rune{11}
var CTRL_L = []rune{12}
var CTRL_N = []rune{14}
var CTRL_O = []rune{15}
var CTRL_P = []rune{16}
var CTRL_Q = []rune{17}
var CTRL_R = []rune{18}
var CTRL_S = []rune{19}
var CTRL_T = []rune{20}
var CTRL_U = []rune{21}
var CTRL_V = []rune{22}
var CTRL_W = []rune{23}
var CTRL_X = []rune{24}
var CTRL_Y = []rune{25}
var CTRL_Z = []rune{26}

var TAB = []rune{9}
var LINE_FEED = []rune{10}
var CARRIAGE_RETURN = []rune{13}
var BACKSPACE = []rune{127}

var SHIFT_TAB = []rune{27, 91, 90}

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

var (
	black = color("\033[1;30m%s\033[0m")
	red = color("\033[1;31m%s\033[0m")
	green = color("\033[1;32m%s\033[0m")
	yellow = color("\033[1;33m%s\033[0m")
	blue = color("\033[1;34m%s\033[0m")
	magenta = color("\033[1;35m%s\033[0m")
	cyan = color("\033[1;36m%s\033[0m")
	lightGray = color("\033[1;37m%s\033[0m")
	defaultColor = color("\033[1;39m%s\033[0m")
	darkGray = color("\033[1;90m%s\033[0m")
	lightRed = color("\033[1;91m%s\033[0m")
	lightGreen = color("\033[1;92m%s\033[0m")
	lightYellow = color("\033[1;93m%s\033[0m")
	lightBlue = color("\033[1;94m%s\033[0m")
	lightMagenta = color("\033[1;95m%s\033[0m")
	lightCyan = color("\033[1;96m%s\033[0m")
	white = color("\033[1;97m%s\033[0m")
)
