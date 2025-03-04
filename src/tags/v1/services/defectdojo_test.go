package ydk

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

// public credentials for demo purpose
var options = DefectDojoServiceOptions{
	Url:      "https://demo.defectdojo.org",
	User:     "admin",
	Password: "1Defectdojo@demo#appsec",
}

func TestDefectDojoServiceAuth(t *testing.T) {
	service := NewDefectDojoService(options)
	if token, err := service.GetToken(options.User, options.Password); err == nil {
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
	service := NewDefectDojoService(options)
	if service.Options.Url != options.Url {
		t.Errorf("Expected URL %s, got %s", options.Url, service.Options.Url)
	}
	if service.Options.User != options.User {
		t.Errorf("Expected user %s, got %s", options.User, service.Options.User)
	}
	if service.Options.Password != options.Password {
		t.Errorf("Expected password %s, got %s", options.Password, service.Options.Password)
	}
	if service.Consumer == nil {
		t.Errorf("Expected non-nil HTTP client, got nil")
	}
}

func TestAuth(t *testing.T) {
	Auth()
}

func Auth() {

	url := "https://demo.defectdojo.org/api/v2/api-token-auth/"
	method := "POST"

	payload := strings.NewReader("username=admin&password=1Defectdojo%40demo%23appsec")

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Basic YWRtaW46MURlZmVjdGRvam9AZGVtbyNhcHBzZWM=")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}
