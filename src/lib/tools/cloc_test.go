package tools_test

import (
	"os"
	"testing"

	// reports "yellowteam/lib/report"
	tools "yellowteam/lib/tools"
)

func TestCountLinesOfCode_FileNotExist(t *testing.T) {
	result := program.CountLinesOfCode("nonexistent_file.go")
	if result != nil {
		t.Errorf("Expected nil, got %v", result)
	}
}

var program = tools.ClocProgram{}

func TestCountLinesOfCode_Directory(t *testing.T) {
	result := program.CountLinesOfCode("/tmp")
	if result == nil {
		t.Errorf("Expected non-nil result, got nil.")
	}
	// t.Log(result.Json(true))
	// report := reports.NewSarif()
	// result.Sarif(report)
	// t.Log(report.Json(true))
	// program.Uninstall()
}

func TestCountLinesOfCode_ValidFile(t *testing.T) {
	file, err := os.CreateTemp("", "testfile*.go")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(file.Name())

	_, err = file.WriteString("package main\n\nfunc main() {}\n")
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	file.Close()

	result := program.CountLinesOfCode(file.Name())
	if result == nil {
		t.Errorf("Expected non-nil result, got nil")
	}
}
