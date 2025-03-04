package ydk

type YdkConsoleColorIndex int
const (
	NoColor = iota
	NoBackground
	Black
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
	Gray
	Purple
	Orange
)
var YdkConsoleColors = map[YdkConsoleColorIndex]string{
	NoColor:      "0",
	NoBackground: "0",
	Black:        "30",
	Red:          "31",
	Green:        "32",
	Yellow:       "33",
	Blue:         "34",
	Magenta:      "35",
	Cyan:         "36",
	White:        "37",
	Gray:         "90",
	Purple:       "95",
	Orange:       "91",
}
var ColorNameToIndex = map[string]YdkConsoleColorIndex{
	"nocolor":      NoColor,
	"nobackground": NoBackground,
	"black":        Black,
	"red":          Red,
	"green":        Green,
	"yellow":       Yellow,
	"blue":         Blue,
	"magenta":      Magenta,
	"cyan":         Cyan,
	"white":        White,
	"gray":         Gray,
	"purple":       Purple,
	"orange":       Orange,
}
type YdkConsoleVariantIndex int

const (
	Normal = iota
	Bold
	Dim
	Italic
	Underline
	Blink
	Reverse
	Hidden
	Foreground
	Background
	Light
	Dark
)
var YdkConsoleColorVariants = map[YdkConsoleVariantIndex]string{
	Normal:     "0",
	Bold:       "1",
	Dim:        "2",
	Italic:     "3",
	Underline:  "4",
	Blink:      "5",
	Reverse:    "7",
	Hidden:     "8",
	Foreground: "38",
	Background: "48",
	Light:      "1",
	Dark:       "2",
}