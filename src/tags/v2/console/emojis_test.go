package ydk

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestListEmojis(t *testing.T) {
	expectedEmojis := map[string]string{
		"smile":   "😄",
		"heart":   "❤️",
		"thumbsup": "👍",
	}

	emojis = expectedEmojis

	result := ListEmojis()
	fmt.Println(result)

	if !reflect.DeepEqual(result, expectedEmojis) {
		t.Errorf("ListEmojis() = %v, want %v", result, expectedEmojis)
	}
}

func TestInterpolateEmojis(t *testing.T) {
	// Test case 1: Text with emojis to interpolate
	text := "I am feeling :smile: today!"
	expectedResult := "I am feeling 😄 today!"
	result := InterpolateEmojis(text)
	if result != expectedResult {
		t.Errorf("InterpolateEmojis(%q) = %q, want %q", text, result, expectedResult)
	}
	fmt.Println(result)

	// Test case 2: Text without emojis to interpolate
	text = "No emojis here!"
	expectedResult = "No emojis here!"
	result = InterpolateEmojis(text)
	if result != expectedResult {
		t.Errorf("InterpolateEmojis(%q) = %q, want %q", text, result, expectedResult)
	}
	fmt.Println(result)

	// Test case 3: Text with multiple emojis to interpolate
	text = "I :heart: :smile: :thumbsup: emojis!"
	expectedResult = "I ❤️ 😄 👍 emojis!"
	result = InterpolateEmojis(text)
	if result != expectedResult {
		t.Errorf("InterpolateEmojis(%q) = %q, want %q", text, result, expectedResult)
	}
	fmt.Println(result)

	text = "I :heart: :smile: :thumbsup: and :random: emojis!"
	expectedResult = "I ❤️ 😄 👍"
	result = InterpolateEmojis(text)
	fmt.Println(result)
	// if result not contains
	if strings.Contains(result, ":random:") {
		t.Errorf("InterpolateEmojis(%q) = %q, want %q", text, result, expectedResult)
	}
}