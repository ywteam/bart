package ydk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

func TestLanguagesScanner(t *testing.T) {
	scanners := Scanners()
	if len(scanners) == 0 {
		t.Errorf("Expected non-empty list of scanners, got empty list")
	}
	scanner := scanners[0]
	if scanner.Id == "" {
		t.Errorf("Expected non-empty scanner ID, got empty string")
	}

	if result, err := scanner.Entrypoint([]string{"/workspace"}); err == nil {
		if result.ExitCode != 0 {
			t.Errorf("Expected exit code 0, got %d", result.ExitCode)
		}
		if result.Output.String() == "" {
			t.Errorf("Expected non-empty output, got empty string")
		}
		data := map[string]interface{}{}
		if err := json.Unmarshal(result.Output.Bytes(), &data); err != nil {
			t.Errorf("Expected valid JSON, got %s", result.Output.String())
		}
		fmt.Println(data)
	} else {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestScannerVersion(t *testing.T) {
	scanners := Scanners()
	if len(scanners) == 0 {
		t.Errorf("Expected non-empty list of scanners, got empty list")
	}
	scanner := scanners[0]
	if scanner.Id == "" {
		t.Errorf("Expected non-empty scanner ID, got empty string")
	}

	if version, err := scanner.Version(); err == nil {
		if version == "" {
			t.Errorf("Expected non-empty version, got empty string")
		}
		if version != "1.96" {
			t.Errorf("Expected version 1.96, got %s", version)
		}

		fmt.Println(version)
	} else {
		t.Errorf("Expected no error, got %v", err)
	}

}

func TestScannerEntrypoint(t *testing.T) {
	scanners := Scanners()
	if len(scanners) == 0 {
		t.Errorf("Expected non-empty list of scanners, got empty list")
	}
	scanner := scanners[0]
	if scanner.Id == "" {
		t.Errorf("Expected non-empty scanner ID, got empty string")
	}

	if result, err := scanner.Entrypoint([]string{"--version"}); err == nil {
		if result.ExitCode != 0 {
			t.Errorf("Expected exit code 0, got %d", result.ExitCode)
		}
		if result.Output.String() == "" {
			t.Errorf("Expected non-empty output, got empty string")
		}
		fmt.Println(result.Output.String())
	} else {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestScanners(t *testing.T) {
	scanners := Scanners()
	if len(scanners) == 0 {
		t.Errorf("Expected non-empty list of scanners, got empty list")
	}
}

func TestYdk(t *testing.T) {
	expected := "cbb46398-a79e-4afe-9672-badabf6075e7"
	Ydk()
	actual := captureOutput(Ydk)
	if actual != expected {
		t.Errorf("Expected %s, but got %s", expected, actual)
	}
}

func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)

	return strings.Trim(buf.String(), "\n")
}
