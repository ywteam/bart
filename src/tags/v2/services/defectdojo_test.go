package ydk

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"testing"
)

// public credentials for demo purpose
var defectDojoOptions = DefectDojoServiceOptions{
	Url:      "https://demo.defectdojo.org",
	User:     "admin",
	Password: "1Defectdojo@demo#appsec",
	Token:    "497288f15b39193d39edc2c21a8795f5f02eb2d7",
}

func TestDefectDojoServiceAuth(t *testing.T) {
	service := NewDefectDojoService(defectDojoOptions)
	if token, err := service.GetToken(defectDojoOptions.User, defectDojoOptions.Password); err == nil {
		if token == "" {
			t.Errorf("Expected non-empty token, got empty string")
		}
	} else {
		t.Errorf("Expected no error, got %v", err)
	}
	if service.Options.Token == "" {
		t.Errorf("Expected non-empty token, got empty string")
	}
	fmt.Println(service.Options.Token)
}
func TestNewDojoService(t *testing.T) {
	service := NewDefectDojoService(defectDojoOptions)
	fmt.Println(service)
	if service.Consumer.Options.Url.String() != defectDojoOptions.Url+"/api/v2" {
		t.Errorf("Expected URL %s, got %s", defectDojoOptions.Url, service.Consumer.Options.Url.String())
	}
	if service.Options.Url != defectDojoOptions.Url {
		t.Errorf("Expected URL %s, got %s", defectDojoOptions.Url, service.Options.Url)
	}
	if service.Options.User != defectDojoOptions.User {
		t.Errorf("Expected user %s, got %s", defectDojoOptions.User, service.Options.User)
	}
	if service.Options.Password != defectDojoOptions.Password {
		t.Errorf("Expected password %s, got %s", defectDojoOptions.Password, service.Options.Password)
	}
	if service.Consumer == nil {
		t.Errorf("Expected non-nil HTTP client, got nil")
	}
	if !service.Options.HasToken() {
		t.Errorf("Expected token, got none")
	}
}

func TestAuth(t *testing.T) {
	url := defectDojoOptions.Url + "/api/v2/api-token-auth/"
	method := "POST"
	payload := map[string]string{
		"username": defectDojoOptions.User,
		"password": defectDojoOptions.Password,
	}
	client := &http.Client{}
	var payloadBuffer bytes.Buffer
	multiPartWriter := multipart.NewWriter(&payloadBuffer)
	for key, value := range payload {
		_ = multiPartWriter.WriteField(key, value)
	}
	multiPartWriter.Close()
	req, err := http.NewRequest(method, url, &payloadBuffer)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	// req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Type", multiPartWriter.FormDataContentType())
	req.Header.Add("Authorization", "Basic YWRtaW46MURlZmVjdGRvam9AZGVtbyNhcHBzZWM=")
	// service := NewDefectDojoService(options)
	// // service.Consumer.DebugRequest(req)
	fmt.Println("[Consumer] Request", req)
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	defer res.Body.Close()
	fmt.Println("[Consumer] Response", res)
	// service.Consumer.DebugResponse(res)
	if res.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", res.StatusCode)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	fmt.Println(string(body))
}

func TestAuth2(t *testing.T) {
	// service := NewDefectDojoService(options)

	url := "https://demo.defectdojo.org/api/v2/api-token-auth/"
	method := "POST"

	payload := strings.NewReader("username=admin%0A&password=1Defectdojo%40demo%23appsec")

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}
	// service.Consumer.DebugRequest(req)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}
	defer res.Body.Close()
	// service.Consumer.DebugResponse(res)
	if res.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", res.StatusCode)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}
	fmt.Println(string(body))
}
