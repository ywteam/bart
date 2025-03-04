package console

import (
	"bytes"
	"os"
	"testing"
)

func TestConsole_Print(t *testing.T) {
	var buf bytes.Buffer
	c := New(&buf, os.Stderr)
	c.Print("Hello, World!")
	if buf.String() != "Hello, World!" {
		t.Errorf("Expected 'Hello, World!', got '%s'", buf.String())
	}
}

func TestConsole_Println(t *testing.T) {
	var buf bytes.Buffer
	c := New(&buf, os.Stderr)
	c.Println("Hello, World!")
	if buf.String() != "Hello, World!\n" {
		t.Errorf("Expected 'Hello, World!\\n', got '%s'", buf.String())
	}
}

func TestConsole_Printf(t *testing.T) {
	var buf bytes.Buffer
	c := New(&buf, os.Stderr)
	c.Printf("Hello, %s!", "World")
	if buf.String() != "Hello, World!" {
		t.Errorf("Expected 'Hello, World!', got '%s'", buf.String())
	}
}

func TestConsole_Error(t *testing.T) {
	var buf bytes.Buffer
	c := New(os.Stdout, &buf)
	c.Error("Error message")
	if buf.String() != "Error message" {
		t.Errorf("Expected 'Error message', got '%s'", buf.String())
	}
}

func TestConsole_Errorln(t *testing.T) {
	var buf bytes.Buffer
	c := New(os.Stdout, &buf)
	c.Errorln("Error message")
	if buf.String() != "Error message\n" {
		t.Errorf("Expected 'Error message\\n', got '%s'", buf.String())
	}
}

func TestConsole_Errorf(t *testing.T) {
	var buf bytes.Buffer
	c := New(os.Stdout, &buf)
	c.Errorf("Error: %s", "message")
	if buf.String() != "Error: message" {
		t.Errorf("Expected 'Error: message', got '%s'", buf.String())
	}
}

func TestConsole_Style(t *testing.T) {
	var buf bytes.Buffer
	c := New(&buf, os.Stderr)
	c.Style("bold", "underline", "reset")
	expected := "\033[1m\033[4m\033[0m"
	if buf.String() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, buf.String())
	}
}

func TestConsole_Colorize(t *testing.T) {
	c := New(os.Stdout, os.Stderr)
	result := c.Colorize("red", "bold", "reset")
	expected := "\033[31m\033[1m\033[0m"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestConsole_Clear(t *testing.T) {
	var buf bytes.Buffer
	c := New(&buf, os.Stderr)
	c.Clear()
	expected := "\033[H\033[2J"
	if buf.String() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, buf.String())
	}
}

func TestConsole_Title(t *testing.T) {
	var buf bytes.Buffer
	c := New(&buf, os.Stderr)
	c.Title("Test Title")
	expected := "\033]0;Test Title\007"
	if buf.String() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, buf.String())
	}
}

func TestConsole_Cursor(t *testing.T) {
	var buf bytes.Buffer
	c := New(&buf, os.Stderr)
	c.Cursor(10, 20)
	expected := "\033[10;20H"
	if buf.String() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, buf.String())
	}
}

func TestConsole_CursorUp(t *testing.T) {
	var buf bytes.Buffer
	c := New(&buf, os.Stderr)
	c.CursorUp(5)
	expected := "\033[5A"
	if buf.String() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, buf.String())
	}
}

func TestConsole_CursorDown(t *testing.T) {
	var buf bytes.Buffer
	c := New(&buf, os.Stderr)
	c.CursorDown(5)
	expected := "\033[5B"
	if buf.String() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, buf.String())
	}
}

// func TestConsole_Prompt(t *testing.T) {
// 	var buf bytes.Buffer
// 	c := New(&buf, os.Stderr)
// 	input := "user input"
// 	buf.WriteString(input + "\n")
// 	result, err := c.Prompt("Enter something: ")
// 	if err != nil {
// 		t.Errorf("Unexpected error: %v", err)
// 	}
// 	if result != input {
// 		t.Errorf("Expected '%s', got '%s'", input, result)
// 	}
// }
