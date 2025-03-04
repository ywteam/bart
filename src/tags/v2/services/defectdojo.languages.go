package ydk

import (
	"fmt"
	"os"
)

type DefectDojoLanguageService struct {
	Api *DefectDojoService
}

func (s *DefectDojoLanguageService) Import(productId int, report string) (*DefectDojoResponse, error) {
	// validate if file exits
	if len(report) == 0 {
		return nil, fmt.Errorf("missing report")
	}
	if fileInfo, err := os.Stat(report); err != nil {
		return nil, err
	} else if fileInfo.IsDir() {
		return nil, fmt.Errorf("report is a directory")
	}
	file, err := os.Open(report)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	resp, err := s.Api.Consumer.Form("/import-languagues/", HttpRequestForm{
		UrlEncoded: map[string]interface{}{
			"file": "@" + report,
			"product": productId,
		},
	})
	if err != nil {
		return nil, err
	}
	var response DefectDojoResponse
	if err := resp.Json(&response); err != nil {
		return nil, err
	}
	return &response, nil
}
func (s *DefectDojoLanguageService) List() (*DefectDojoLanguageResponse, error) {
	resp, err := s.Api.Consumer.Get("/languages/", nil)
	if err != nil {
		return nil, err
	}
	var languages DefectDojoLanguageResponse
	if err := resp.Json(&languages); err != nil {
		return nil, err
	}
	return &languages, nil
}

type DefectDojoLanguageModel struct {
	Id         int    `json:"id"`
	Files      int    `json:"files"`
	Blank      int    `json:"blank"`
	Comment    int    `json:"comment"`
	Code       int    `json:"code"`
	Created    string `json:"created"`
	LanguageId int    `json:"language"`
	ProductId  int    `json:"product"`
	UserId     int    `json:"user"`
}
type DefectDojoLanguageResponse struct {
	DefectDojoResponse
	Results []DefectDojoLanguageModel `json:"results"`
}
