package console

import (
	"fmt"
	"os"
	"io"
	"strconv"
	"sync"
	"strings"
)
var Colors = map[string]string{
	"black": "30",
	"red": "31",
	"green": "32",
	"yellow": "33",
	"blue": "34",
	"magenta": "35",
	"cyan": "36",
	"white": "37",
	"gray": "90",
	"purple": "95",
	"orange": "91",	
}
var Styles = map[string]string{
	// "reset": "\033[0m",
	"reset": "0",
	"bold": "1",
	"dim": "2",
	"italic": "3",
	"underline": "4",
	"blink": "5",
	"reverse": "7",
	"hidden": "8",	
}

type Console struct {
	stdout io.Writer	//*os.File
	stderr io.Writer	//*os.File	
}
func (c *Console) Style(styles ...string) *Console {
	for _, style := range styles {
		if _, ok := Styles[style]; !ok {
			// fmt.Print(style)
			fmt.Fprint(c.stdout, style)
			 continue
		}
		fmt.Fprintf(c.stdout, "\033[%sm", Styles[style])
		// fmt.Printf("\033[%sm", Styles[style])
	}
	return c
}
func (c *Console) Colorize(styles ...string) string{
	var result = []string{}
	for _, style := range styles {
		if _, ok := Styles[style]; !ok {
			result = append(result, style)
			continue
		}
		result = append(result, fmt.Sprintf("\033[%sm", Styles[style]))
	}
	return strings.Join(result, "")
}
func (c *Console) Clear() *Console {
	fmt.Fprint(c.stdout, "\033[H\033[2J")
	return c
}
func (c *Console) Title(text string) *Console {
	fmt.Fprintf(c.stdout, "\033]0;%s\007", text)
	return c
}
func (c *Console) Cursor(x int, y int) *Console {
	fmt.Fprintf(c.stdout, "\033[%d;%dH", x, y)
	return c
}
func (c *Console) CursorUp(n int) *Console {
	fmt.Fprintf(c.stdout, "\033[%dA", n)
	return c
}
func (c *Console) CursorDown(n int) *Console {
	fmt.Fprintf(c.stdout, "\033[%dB", n)
	return c
}
// func (c *Console) Prompt(message string) (string, error) {
// 	fmt.Fprint(os.Stdin, message)
// 	var input string
// 	_, err := fmt.Fscanln(os.Stdin, &input)
// 	return input, err	
// }
func (c *Console) Print(a ...any) *Console {
	if _, ok := fmt.Fprint(c.stdout, a...); ok != nil {
		c.Errorln(ok)
	}
	return c
	// if _, ok := fmt.Print(a...); ok != nil {
	// 	c.Errorln(ok)
	// }
	// return c
}
func (c *Console) Printf(format string, a ...any) *Console {	
	if _, ok := fmt.Fprintf(c.stdout, format, a...); ok != nil {
		c.Errorln(ok)
	}
	return c
}
func (c *Console) Println(a ...any) *Console {
	if _, ok := fmt.Fprintln(c.stdout, a...); ok != nil {
		c.Errorln(ok)
	}
	return c
}
func (c *Console) Error(a ...any) *Console {
	if _, ok := fmt.Fprint(c.stderr, a...); ok != nil {
		c.Errorln(ok)
	}
	return c
}
func (c *Console) Errorf(format string, a ...any) *Console {
	if _, ok := fmt.Fprintf(c.stderr, format, a...); ok != nil {
		c.Errorln(ok)
	}
	return c
}
func (c *Console) Errorln(a ...any) *Console {
	if _, ok := fmt.Fprintln(c.stderr, a...); ok != nil {
		c.Fatal(ok)
	}
	return c
}
func (c *Console) Write(p []byte) *Console {
	if _, ok := c.stdout.Write(p); ok != nil {
		c.Fatal(ok)
	}
	return c
}
func (c *Console) Writeln(p []byte) *Console {
	if _, ok := c.stdout.Write(append(p, '\n')); ok != nil {
		c.Fatal(ok)
	}
	return c
}
func (c *Console) Fatal(a ...any) {
	fmt.Fprint(c.stderr, a...)
	c.Exit(1)
}
func (c *Console) Fatalf(format string, a ...any) {
	fmt.Fprintf(c.stderr, format, a...)
	c.Exit(1)
}
func (c *Console) Fatalln(a ...any) {
	fmt.Fprintln(c.stderr, a...)
	c.Exit(1)
}
func (c *Console) Exit(code int) {
	os.Exit(code)
}


var instance *Console
var once sync.Once
func New(stdout, stderr io.Writer) *Console {
	return &Console{
		stdout: stdout,
		stderr: stderr,
	}
}
func Default() *Console {
	return instance
}

func init(){
	for color, code := range Colors {		
		Styles[color] = code
		codeInt, _ := strconv.Atoi(code)
		Styles[color + ":bg"] = fmt.Sprintf("%d", codeInt + 10)
		Styles[color + ":light"] = fmt.Sprintf("%d", codeInt + 60)
		Styles[color + ":light:bg"] = fmt.Sprintf("%d", codeInt + 70)
	}
	once.Do(func(){
		instance = New(os.Stdout, os.Stderr)
	})
}