package ydk

import (
	"fmt"
	"testing"
)

func TestDefectDojoLanguageService_List(t *testing.T) {
	service := NewDefectDojoService(defectDojoOptions)
	languages, err := service.Languages.List()
	if err != nil {
		t.Error(err)
	}
	if languages.Count == 0 {
		t.Error("Expected languages, got none")
	}
	fmt.Println("TestDefectDojoLanguageService_List()", languages)
}