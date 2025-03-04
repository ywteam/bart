package ydk

import "testing"
func TestHttpConsumerAuth(t *testing.T) {
	svc := NewHttpExampleService();
	if svc.Consumer.Http == nil {
		t.Errorf("Expected non-nil HTTP client, got nil")
	}
	if svc.Consumer.BaseUrl != "https://demo.defectdojo.org/api/v2" {
		t.Errorf("Expected base URL %s, got %s", "https://demo.defectdojo.org/api/v2", svc.Consumer.BaseUrl)
	}
	example, err := svc.GetExample()
	if err != nil {
		t.Errorf("Error getting example: %s", err)
	}
	if example.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %d", example.StatusCode)
	}	
	
}func TestHttpHttpConsumerService_Curl(t *testing.T) {
	service := NewHttpHttpConsumerService(HTTP_CONSUMER_DEFAULTS)
	req, err := http.NewRequest("GET", "https://example.com", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	resp, err := service.Curl(req)
	if err != nil {
		t.Fatalf("Failed to execute request: %v", err)
	}

	// Add assertions here to validate the response
}