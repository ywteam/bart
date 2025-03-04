package osv_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/ywteam/ydk-go/apis/osv"
)

func TestQueryBatchAPI(t *testing.T) {
	// Define o mock response da API
	mockResponse := osv.QueryBatchResponse{
		Results: []osv.QueryResult{
			{
				Vulnerabilities: []osv.Vulnerability{
					{
						ID:       "CVE-2024-12345",
						Modified: "2024-11-01T00:00:00Z",
						Details:  "Example vulnerability details",
					},
				},
			},
		},
	}

	// Serializa o mock response para JSON
	mockResponseData, err := json.Marshal(mockResponse)
	if err != nil {
		t.Fatalf("failed to marshal mock response: %v", err)
	}

	// Configura o servidor HTTP de teste
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verifica o método HTTP
		if r.Method != http.MethodPost {
			t.Errorf("expected POST method, got %s", r.Method)
		}

		// Verifica o header
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
		}

		// Retorna o mock response
		w.WriteHeader(http.StatusOK)
		w.Write(mockResponseData)
	}))
	defer server.Close()

	// Define o payload de teste
	testRequest := osv.QueryBatchRequest{
		Queries: []osv.Query{
			{
				Package: &osv.PackageQuery{
					Name:      "example-package",
					Ecosystem: "npm",
				},
			},
		},
	}

	// Chama a função com o servidor de teste
	response, err := osv.QueryBatchAPI(server.URL, testRequest)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verifica a resposta
	if len(response.Results) != 1 {
		t.Errorf("expected 1 result, got %d", len(response.Results))
	}

	vulns := response.Results[0].Vulnerabilities
	if len(vulns) != 1 {
		t.Errorf("expected 1 vulnerability, got %d", len(vulns))
	}

	vuln := vulns[0]
	if vuln.ID != "CVE-2024-12345" {
		t.Errorf("expected vulnerability ID 'CVE-2024-12345', got '%s'", vuln.ID)
	}
	if vuln.Details != "Example vulnerability details" {
		t.Errorf("expected details 'Example vulnerability details', got '%s'", vuln.Details)
	}
	if vuln.Modified != "2024-11-01T00:00:00Z" {
		t.Errorf("expected modified '2024-11-01T00:00:00Z', got '%s'", vuln.Modified)
	}
}
func TestQueryBatchAPIIntegration(t *testing.T) {
	// URL da API pública
	apiURL := "https://api.osv.dev/v1/querybatch"

	// Payload real
	request := osv.QueryBatchRequest{
		Queries: []osv.Query{
			{
				Package: &osv.PackageQuery{
					Name:      "express",
					Ecosystem: "npm",
				},
			},			
			{
				Commit: "d2d2b417f9b64b6785b8c3ad1ef6027b02f122c3", // Commit do exemplo
			},
		},
	}

	// Chama a API real
	response, err := osv.QueryBatchAPI(apiURL, request)
	if err != nil {
		t.Fatalf("error querying OSV API: %v", err)
	}

	// Verifica se há resultados na resposta
	if len(response.Results) == 0 {
		t.Fatalf("expected results but got none")
	}

	// Verifica vulnerabilidades em um dos resultados
	for i, result := range response.Results {
		t.Logf("Result %d:", i+1)
		if len(result.Vulnerabilities) == 0 {
			t.Logf("No vulnerabilities found for result %d", i+1)
		} else {
			for _, vuln := range result.Vulnerabilities {
				t.Logf("Vulnerability ID: %s", vuln.ID)
				t.Logf("Details: %s", vuln.Details)
				t.Logf("Modified: %s", vuln.Modified)
			}
		}
	}
}
func TestQueryBatchAPIPurlAndCommit(t *testing.T) {
	// URL da API pública
	apiURL := "https://api.osv.dev/v1/querybatch"

	// Payload do teste
	request := osv.QueryBatchRequest{
		Queries: []osv.Query{
			{
				Package: &osv.PackageQuery{
					PURL: "pkg:pypi/mlflow@0.4.0",
				},
			},
			{
				Commit: "6879efc2c1596d11a6a6ad296f80063b558d5e0f",
			},
			{
				Package: &osv.PackageQuery{
					Ecosystem: "PyPI",
					Name:      "jinja2",
				},
				Version: "2.4.1",
			},
		},
	}

	// Faz a requisição para a API real
	response, err := osv.QueryBatchAPI(apiURL, request)
	if err != nil {
		t.Fatalf("error querying OSV API: %v", err)
	}

	// Verifica se a resposta possui resultados
	if len(response.Results) == 0 {
		t.Fatalf("expected results but got none")
	}

	// Processa os resultados
	for i, result := range response.Results {
		t.Logf("Result %d:", i+1)
		if len(result.Vulnerabilities) == 0 {
			t.Logf("No vulnerabilities found for result %d", i+1)
		} else {
			for _, vuln := range result.Vulnerabilities {
				t.Logf("Vulnerability ID: %s", vuln.ID)
				t.Logf("Details: %s", vuln.Details)
				t.Logf("Modified: %s", vuln.Modified)
			}
		}
	}
}
