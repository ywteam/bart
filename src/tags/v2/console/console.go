package ydk

import (
	"fmt"
	"strings"
)

func Colorize(text string, color YdkConsoleColorIndex, variant YdkConsoleVariantIndex) string {
	return fmt.Sprintf("\033[%s;%sm%s\033[0m", YdkConsoleColorVariants[variant], YdkConsoleColors[color], text)
}
func ColorIndex(colorName string) YdkConsoleColorIndex {
	if colorName == "" {
		return NoColor
	}
	colorName = strings.ToLower(colorName)
	if index, exists := ColorNameToIndex[colorName]; exists {
		return index
	}
	return NoColor
}
func Prompt(message string) (string, error) {
	fmt.Print(message)
	var input string
	_, err := fmt.Scanln(&input)
	return input, err
}
func Table(headers []string, data [][]string) {	
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
		}
		fmt.Println()
	}
}