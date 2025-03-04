package defectdojo_test

import (
	"strings"
	"testing"

	defectdojo "github.com/ywteam/ydk-go/apis/defectdojo"
)

// public credentials for demo purpose
// https://github.com/DefectDojo/django-DefectDojo#demo
var defectDojoOptions = defectdojo.DefectDojoServiceOptions{
	Url: "https://demo.defectdojo.org",
	Auth: defectdojo.DefectDojoServiceBasicAuth{
		Username: "admin",
		Password: "1Defectdojo@demo#appsec",
	},
}
func TestNewDojoService(t *testing.T) {
	service, err := defectdojo.NewDefectDojoService(defectDojoOptions)
	if err != nil {
		t.Fatalf("Failed to create DefectDojoService: %v", err)
	}
	if strings.HasPrefix(defectDojoOptions.Url, service.Client.Options.BaseUrl.String()) {
		t.Errorf("Expected URL %s, got %s", defectDojoOptions.Url, service.Client.Options.BaseUrl.String())
	}
	if service.Options.Url != defectDojoOptions.Url {
		t.Errorf("Expected URL %s, got %s", defectDojoOptions.Url, service.Options.Url)
	}
	if service.Options.Auth == nil {
		t.Errorf("Expected non-nil Auth, got nil")
	}
}
func TestDefectDojoServiceAuth(t *testing.T) {
	service, err := defectdojo.NewDefectDojoService(defectDojoOptions)
	if err != nil {
		t.Fatalf("Failed to create DefectDojoService: %v", err)
	}
	if token, err := service.Auth(defectDojoOptions.Auth); err == nil {
		if token == "" {
			t.Errorf("Expected non-empty token, got empty string")
		}
		t.Logf("Token: %s", token)
	} else {
		t.Errorf("Expected no error, got %v", err)
	}
	if service.Options.Auth == nil {
		t.Errorf("Expected non-empty token, got empty string")
	}
}

