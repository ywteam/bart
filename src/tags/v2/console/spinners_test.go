package ydk

import (
	"math/rand"
	"testing"
	"time"
)

func TestListSpinners(t *testing.T) {
	spinners := ListSpinners()
	if len(spinners) == 0 {
		t.Errorf("Expected non-empty map of spinners, got empty map")
	}
	// Add more assertions if needed
}
func TestRandomSpinner(t *testing.T) {
	key, value := RandomSpinner()
	if key == "0" {
		t.Errorf("Expected non-empty key, got %s", key)
	}
	if value.Interval == 0 {
		t.Errorf("Expected non-zero interval, got %f", value.Interval)
	}
	if len(value.Frames) == 0 {
		t.Errorf("Expected non-empty frames, got empty frames")
	}
	// Add more assertions if needed
}

func TestStart(t *testing.T) {
	// Test case 1: Valid spinner key
	key := "dots"
	text := "Test message"
	randomTime := 1 + time.Duration(rand.Intn(2))
	operation := func(message *string) bool {
		
		time.Sleep(randomTime * time.Second)
		*message = text + ". Still working, please wait..."
		randomTime = 1 + time.Duration(rand.Intn(3)) - 1
		time.Sleep(randomTime * time.Second)
		SpinStop()
		return true //randomTime%2 == 0
	}
	result := SpinStart(key, text, operation)

    // Assert that the result is true
    if result != true {
        t.Errorf("Expected true, but got %v", result)
    }
}