package main

import (
	"encoding/json"
	"os"
	"testing"
)

func WriteFile(file, content string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(content)
	return err
}
func TestOpaService_Test(t *testing.T){
	opa := &OpaService{}
	policyDir := os.TempDir()
	err := WriteFile(policyDir + "/zerotrust.rego", opa.ZeroTrustPolicy())
	if err != nil {
		t.Errorf("Failed to write policy: %v", err)
	}
	err = WriteFile(policyDir + "/public.rego", opa.PublicPolicy())
	if err != nil {
		t.Errorf("Failed to write policy: %v", err)
	}
	result := Bash(`
		docker stop opa &> /dev/null
		docker rm opa &> /dev/null
		docker run --name opa -d \
			-u root \
			-v `+policyDir+ `:/policies \
			-p 8181:8181 openpolicyagent/opa run --server --addr :8181

		until curl -s http://localhost:8181/health &> /dev/null; do
			echo "Waiting for OPA server to start"
			sleep 1
		done
		echo "OPA server started"
		docker exec opa opa build -b /policies -t rego 
		# docker exec opa apk add nodejs npm
		# docker exec opa npm install -g k6
		docker stop opa &> /dev/null
		docker rm opa &> /dev/null
	`)
	for line := range result.Chan {
		t.Logf("Line: %s", line)
	}
	t.Logf("Result: %v", result)
	t.Logf("Output: %s", result.String())
	if result.HasError() {
		t.Errorf("Failed to test OPA service: %v", result.Error)
	}
	// jsonResult, err := json.MarshalIndent(result, "", "  ")
	// if err != nil {
	// 	t.Errorf("Failed to marshal result: %v", err)
	// } else {
	// 	t.Logf("Result: %s", jsonResult)
	// }
}
func TestOpaService_Start(t *testing.T) {
	opa := &OpaService{}
	result := opa.Start()
	if result.HasError() {
		t.Errorf("Failed to start OPA service: %v", result.Error)
	}
}

func TestOpaService_Stop(t *testing.T) {
	opa := &OpaService{}
	result := opa.Stop()
	if result.HasError() {
		t.Errorf("Failed to stop OPA service: %v", result.Error)
	}
}
func Test_Stream(t *testing.T) {
	result := Bash(`
		for i in $(seq 1 140); do
			echo "Line $i"
			sleep 0.2
		done		
	`, "")
	for line := range result.Chan {
		t.Logf("Line: %s", line)
	}
	t.Logf("Result: %v", result)
	t.Logf("Output: %s", result.String())
}
func TestBash(t *testing.T) {
	result := Bash("echo 'Hello, World!'")
	if result.HasError() {
		t.Errorf("Bash command failed: %v", result.Error)
	}
	expected := "Hello, World!"
	if result.String() != expected {
		t.Log(json.MarshalIndent(result, "", "  "))
		t.Errorf("Expected %s but got %s", expected, result.String())
	}
}

func TestCommand(t *testing.T) {
	result := Command("echo", "Hello, World!")
	if result.HasError() {
		t.Errorf("Command failed: %v", result.Error)
	}
	expected := "Hello, World!"
	if result.String() != expected {
		t.Errorf("Expected %s but got %s", expected, result.String())
	}
}