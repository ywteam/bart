package fetch_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"yellowteam/lib/fetch"

	"github.com/stretchr/testify/assert"
)

func TestFetchClient_Methods(t *testing.T) {
	methods := []string{"Get", "Post", "Put", "Patch", "Delete", "Head", "Options"}
	for _, method := range methods {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, strings.ToUpper(method), r.Method)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message": "ok"}`))
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("X-Custom-Header", "custom")
		}))
		defer ts.Close()
		// Parse the test server URL
		baseURL, err := url.Parse(ts.URL)
		assert.NoError(t, err)
		// Create a new FetchClient with the test server URL
		client := fetch.NewClient(&fetch.Options{
			BaseUrl: *baseURL,
			Http: &http.Client{
				Timeout: 10 * time.Second,
			},
		})
		// Perform the request
		if method == "Get" {
			resp, err := client.Get("/test")
			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, http.StatusOK, resp.Code)
			assert.Equal(t, `{"message": "ok"}`, resp.Text())
		} else if method == "Post" {
			resp, err := client.Post("/test", nil)
			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, http.StatusOK, resp.Code)
			assert.Equal(t, `{"message": "ok"}`, resp.Text())
		} else if method == "Put" {
			resp, err := client.Put("/test")
			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, http.StatusOK, resp.Code)
			assert.Equal(t, `{"message": "ok"}`, resp.Text())
		} else if method == "Patch" {
			resp, err := client.Patch("/test")
			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, http.StatusOK, resp.Code)
			assert.Equal(t, `{"message": "ok"}`, resp.Text())
		} else if method == "Delete" {
			resp, err := client.Delete("/test")
			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, http.StatusOK, resp.Code)
			assert.Equal(t, `{"message": "ok"}`, resp.Text())
		} else if method == "Head" {
			resp, err := client.Head("/test")
			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, http.StatusOK, resp.Code)
			assert.Equal(t, ``, resp.Text())
		}

	}
}
func TestFetchClient_HttpBin(t *testing.T) {
	baseURL, err := url.Parse("https://httpbin.org")
	assert.NoError(t, err)

	client := fetch.NewClient(&fetch.Options{
		BaseUrl: *baseURL,
		Http: &http.Client{
			Timeout: 10 * time.Second,
		},
	})

	resp, err := client.Get("/get")
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Text(), `"url": "https://httpbin.org/get"`)
}
func TestFetchClient_Patch(t *testing.T) {
	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PATCH", r.Method)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "patched"}`))
	}))
	defer ts.Close()

	// Parse the test server URL
	baseURL, err := url.Parse(ts.URL)
	assert.NoError(t, err)

	// Create a new FetchClient with the test server URL
	client := fetch.NewClient(&fetch.Options{
		BaseUrl: *baseURL,
		Http: &http.Client{
			Timeout: 10 * time.Second,
		},
	})

	// Perform the PATCH request
	resp, err := client.Patch("/test")
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, `{"message": "patched"}`, resp.Text())
}
