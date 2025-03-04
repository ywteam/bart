package ydk

import (
	"fmt"
	"log"
	"os"
	"strings"
)

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

type YdkConsole struct {
	Colors map[YdkConsoleColorIndex]string
}

func (c *YdkConsole) Colorize(text string, color YdkConsoleColorIndex, variant YdkConsoleVariantIndex) string {
	return fmt.Sprintf("\033[%s;%sm%s\033[0m", YdkConsoleColorVariants[variant], YdkConsoleColors[color], text)
}

func Colorize(text string, color YdkConsoleColorIndex, variant YdkConsoleVariantIndex) string {
	return fmt.Sprintf("\033[%s;%sm%s\033[0m", YdkConsoleColorVariants[variant], YdkConsoleColors[color], text)
}
func (c *YdkConsole) ColorIndex(colorName string) YdkConsoleColorIndex {
	if colorName == "" {
		return NoColor
	}
	colorName = strings.ToLower(colorName)
	if index, exists := ColorNameToIndex[colorName]; exists {
		return index
	}
	return NoColor
}
func (c *YdkConsole) Log(text string) {
	log.Println(text)
}
func (c *YdkConsole) Print(text string) {
	fmt.Println(text)
}
func (c *YdkConsole) Error(text string) {
	fmt.Fprintln(os.Stderr, text)
}
func (c *YdkConsole) Fatal(text string) {
	log.Fatal(text)
}
func (c *YdkConsole) Exit(code int) {
	os.Exit(code)
}
func (c *YdkConsole) Clear() {
	fmt.Print("\033[H\033[2J")
}
func (c *YdkConsole) Title(text string) {
	fmt.Printf("\033]0;%s\007", text)
}
func (c *YdkConsole) Cursor(x int, y int) {
	fmt.Printf("\033[%d;%dH", x, y)
}
func (c *YdkConsole) CursorUp(n int) {
	fmt.Printf("\033[%dA", n)
}
func (c *YdkConsole) CursorDown(n int) {
	fmt.Printf("\033[%dB", n)
}
func Prompt(message string) (string, error) {
	fmt.Print(message)
	var input string
	_, err := fmt.Scanln(&input)
	return input, err

}
func PrintTable(headers []string, data [][]string) {
	
	columnWidths := make([]int, len(headers))
	for i, header := range headers {
		columnWidths[i] = len(header)
		for _, row := range data {
			if len(row[i]) > columnWidths[i] {
				columnWidths[i] = len(row[i])
			}
		}
	}
	border := Colorize("|", Yellow, Normal)
	borderH := Colorize("-", Yellow, Normal)

	formats := make([]string, len(headers))
	for i, width := range columnWidths {
		formats[i] = fmt.Sprintf("%%-%ds ", width)
		// formats[i] = Colorize(formats[i], Cyan, Normal)
	}
	for i, header := range headers {
		fmt.Printf(formats[i]+" " + border + " ", header)
	}
	fmt.Println()
	for _, width := range columnWidths {
		fmt.Print(strings.Repeat(borderH, width+1) + " " + border + " ")
	}
	fmt.Println()
	for _, row := range data {
		for i, col := range row {
			fmt.Printf(formats[i]+" "+border+" ", col)
			// if ! strings.Contains(col, "\n") {
			// 	fmt.Printf(formats[i]+ " | ", col)
			// 	continue
			// }
			// lines := strings.Split(col, "\n")
			// for j, line := range lines {
			// 	if j == 0 {
			// 		fmt.Printf(formats[i]+ " | ", line)
			// 	} else {
			// 		fmt.Printf(formats[i]+ " | ", "")
			// 	}
			// 	fmt.Print(" | ")
			// }

			// // expect \n in the string
			// if strings.Contains(col, "\n") {
			// 	lines := strings.Split(col, "\n")
			// 	for j, line := range lines {
			// 		if j == 0 {
			// 			fmt.Printf(formats[i], line)
			// 		} else {
			// 			fmt.Printf(formats[i], "")
			// 		}
			// 		fmt.Print(" | ")
			// 	}
			// } else {
			// 	fmt.Printf(formats[i], col)
			// 	fmt.Print(" | ")
			// }
		}
		fmt.Println()
	}
}
