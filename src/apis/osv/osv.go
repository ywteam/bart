package osv

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func QueryBatchAPI(apiURL string, request QueryBatchRequest) (*QueryBatchResponse, error) {
	// Serializa o payload para JSON
	payload, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Log the request payload for debugging
	fmt.Printf("Request payload: %s\n", payload)

	// Cria a requisição HTTP POST
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Envia a requisição
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Lê a resposta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Log the response body for debugging
	fmt.Printf("Response body: %s\n", body)

	// Verifica o status da resposta
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %s", body)
	}

	// Desserializa a resposta JSON
	var response QueryBatchResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}
