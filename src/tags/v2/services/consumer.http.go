package ydk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// HttpConsumerOptions represents the options for an HTTP consumer.
type HttpConsumerOptions struct {
	Ctx     context.Context   // The context for the HTTP consumer.
	BaseUrl string            // The base URL for the HTTP consumer.
	Url     url.URL           // The URL for the HTTP consumer.
	Http    *http.Client      // The HTTP client for the HTTP consumer.
	Headers map[string]string // The headers for the HTTP consumer.
	Cookies []http.Cookie     // The cookies for the HTTP consumer.
	Params  map[string]string // The parameters for the HTTP consumer.
	Timeout time.Duration     // The timeout for the HTTP consumer.
}
const (
	HTTP_CONSUMER_DEFAULT_API_KEY = "LSj!U0z^ps&Z!78Tak@f4F*#pVU8$8&0$RH!f4d#x9MnF63T5Zpz5PB6yM4M"
)
// HTTP_CONSUMER_DEFAULTS represents the default HTTP consumer options.
var HTTP_CONSUMER_DEFAULTS = &HttpConsumerOptions{
	Ctx:  context.Background(),
	Http: http.DefaultClient,
	Headers: map[string]string{
		"User-Agent": "Mozilla/5.0",
		"Accept":     "application/json",
		"Y-Version":  "1.0",
		"Y-Api-Key": HTTP_CONSUMER_DEFAULT_API_KEY,
	},
	Cookies: []http.Cookie{},
	Params:  map[string]string{},
	Timeout: time.Duration(10 * time.Second),
}

// Validate validates the HttpConsumerOptions and sets default values if necessary.
func (o *HttpConsumerOptions) Validate() error {
	if o.Ctx == nil {
		o.Ctx = context.Background()
	}
	if o.BaseUrl == "" {
		return fmt.Errorf("missing base URL")
	}
	if o.Headers == nil {
		o.Headers = map[string]string{}	
	}
	if o.Cookies == nil {
		o.Cookies = []http.Cookie{}	
	}
	if o.Params == nil {
		o.Params = map[string]string{}		
	}
	url, err := url.Parse(o.BaseUrl)
	if err != nil {
		return err
	}
	fmt.Println("URL", url, "BaseURL", o.BaseUrl)
	o.Url = *url
	if o.Http == nil {
		o.Http = &http.Client{
			Transport: http.DefaultTransport,
			Timeout:   time.Duration(HTTP_CONSUMER_DEFAULTS.Timeout),
		}
	}
	return nil
}
func (o *HttpConsumerOptions) Path(path string) (*url.URL, error) {
	return url.Parse(o.Url.String() + path)
}

// HttpRequestForm represents the structure of an HTTP request form.
type HttpRequestForm struct {
	Data       interface{} // Data represents the data to be sent in the request.
	UrlEncoded interface{} // UrlEncoded represents the URL-encoded data to be sent in the request.
	RawText    string      // RawText represents the raw text data to be sent in the request.
	RawJson    string      // RawJson represents the raw JSON data to be sent in the request.
	RawXml     string      // RawXml represents the raw XML data to be sent in the request.
	RawHtml    string      // RawHtml represents the raw HTML data to be sent in the request.
	RawBinary  []byte      // RawBinary represents the raw binary data to be sent in the request.
}

// HttpRequest represents an HTTP request.
type HttpRequest struct {
	Method  string            // The HTTP method of the request (e.g., GET, POST, PUT, DELETE).
	Path    string            // The path of the request URL.
	Headers [][]string        // The headers of the request.
	Params  map[string]string // The query parameters of the request.
	Cookies []http.Cookie     // The cookies of the request.
	Form    HttpRequestForm   // The form data of the request.
}

// Validate checks if the HttpRequest object is valid.
// It returns an error if any required fields are missing.
func (r *HttpRequest) Validate() error {
	if r.Method == "" {
		return fmt.Errorf("missing method")
	}
	if r.Path == "" {
		return fmt.Errorf("missing path")
	}
	if r.Headers == nil {
		r.Headers = [][]string{}
	}
	if r.Params == nil {
		r.Params = map[string]string{}
	}
	if r.Cookies == nil {
		r.Cookies = []http.Cookie{}
	}
	return nil
}

// SetHeader sets the specified header key-value pair in the HttpRequest.
func (r *HttpRequest) SetHeader(key string, value string) {
	r.Headers = append(r.Headers, []string{key, value})
}
func (r *HttpRequest) SetContentType(contentType string) {
	// remove existing content-type header
	for i, header := range r.Headers {
		if header[0] == "Content-Type" {
			r.Headers = append(r.Headers[:i], r.Headers[i+1:]...)
			r.SetHeader("Content-Type", contentType)
			return
		}
	}
	r.SetHeader("Content-Type", contentType)
}

// HttResponse represents an HTTP response.
type HttResponse struct {
	Url     string        // The URL of the HTTP request.
	Code    int           // The HTTP status code of the response.
	Status  string        // The HTTP status message of the response.
	Success bool          // Indicates whether the request was successful or not.
	Headers [][]string    // The headers of the HTTP response.
	Cookies []http.Cookie // The cookies received in the response.
	Body    bytes.Buffer  // The body of the HTTP response.
}

// String returns a string representation of the HttpResponse.
// It includes the URL, response code, status, and body of the response.
func (r *HttResponse) String() string {
	return fmt.Sprintf("Response: %s %d %s %s", r.Url, r.Code, r.Status, r.Body.String())
}

// Text returns the response body as a string.
func (r *HttResponse) Text() string {
	return r.Body.String()
}

// Json decodes the HTTP response body into the provided interface.
// It uses the json.NewDecoder function to decode the response body.
// The decoded data is then stored in the provided interface v.
// If there is an error during decoding, it returns the error.
func (r *HttResponse) Json(v interface{}) error {
	decoder := json.NewDecoder(&r.Body)
	return decoder.Decode(&v)
}

// Header returns the value of the specified header key and a boolean indicating
// whether the header key exists in the HTTP response.
func (r *HttResponse) Header(key string) (string, bool) {
	for _, header := range r.Headers {
		if header[0] == key {
			return header[1], true
		}
	}
	return "", false
}

// IHttpConsumer is an interface for making HTTP requests.
type IHttpConsumer interface {
	Curl(req *http.Request) (string, error)
}

// HttpConsumerService represents a service that consumes HTTP requests and handles responses.
type HttpConsumerService struct {
	Options        *HttpConsumerOptions                                // Options contains the configuration options for the HTTP consumer service.
	AuthHandler    func(*http.Request) (*http.Request, error)          // AuthHandler is a function that handles authentication for incoming requests.
	ResponseHander func(*http.Response, error) (*http.Response, error) // ResponseHandler is a function that handles the response from the HTTP request.
}

// HandleAuth handles the authentication for the HTTP consumer service.
// It takes an HTTP request as input and returns a modified request with authentication information, if applicable.
// If an authentication handler is set, it delegates the authentication process to the handler.
// If no authentication handler is set, it simply returns the original request.
// It returns the modified request and any error encountered during the authentication process.
func (c *HttpConsumerService) HandleAuth(req *http.Request) (*http.Request, error) {
	if c.AuthHandler != nil {
		return c.AuthHandler(req)
	}
	return req, nil
}

// HandleResponse handles the HTTP response and returns the response and error.
// If an error occurred during the HTTP request, it returns the response and the error.
// If the HTTP status code is not within the range of 200-399, it returns the response with an error message.
// If a custom response handler is set, it calls the handler function and returns the response and error returned by the handler.
func (c *HttpConsumerService) HandleResponse(res *http.Response, err error) (*http.Response, error) {
	if err != nil {
		return res, err
	}
	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		if c.ResponseHander != nil {
			return c.ResponseHander(res, fmt.Errorf("HTTP Error: %s", res.Status))
		}
		return res, fmt.Errorf("HTTP Error: %s", res.Status)
	}
	if c.ResponseHander != nil {
		return c.ResponseHander(res, err)
	}
	return res, err
}
func (s *HttpConsumerService) DebugResponse(res *http.Response) {
	fmt.Println("[Consumer] HttpConsumerService Debug")
	fmt.Println("[Consumer] Status", res.Status)
	fmt.Println("[Consumer] StatusCode", res.StatusCode)
	fmt.Println("[Consumer] Proto", res.Proto)
	fmt.Println("[Consumer] ProtoMajor", res.ProtoMajor)
	fmt.Println("[Consumer] ProtoMinor", res.ProtoMinor)
	fmt.Println("[Consumer] Header", res.Header)
	if res.Body != nil {
		fmt.Println("[Consumer] Body exists")
		bodyLen := res.ContentLength
		if bodyLen > 0 {
			fmt.Println("[Consumer] Body length", bodyLen)
		} else {
			fmt.Println("[Consumer] Body length unknown")
		}
		body, err := io.ReadAll(res.Body)
		if err == nil {
			fmt.Println("[Consumer] Body", string(body))
		} else {
			fmt.Println("[Consumer] No body", err)
		}
	}
	fmt.Println("[Consumer] Cookies", res.Cookies())
}

// DebugRequest prints debug information about the HTTP request.
func (s *HttpConsumerService) DebugRequest(req *http.Request) {
	fmt.Println("[Consumer] HttpConsumerService Debug")
	fmt.Println("[Consumer] Method", req.Method)
	fmt.Println("[Consumer] URL", req.URL)
	fmt.Println("[Consumer] Header", req.Header)
	if req.Body != nil {
		fmt.Println("[Consumer] Body exists")
		bodyLen := req.ContentLength
		if bodyLen > 0 {
			fmt.Println("[Consumer] Body length", bodyLen)
		} else {
			fmt.Println("[Consumer] Body length unknown")
		}
		clone := req.Clone(req.Context())
		body, err := io.ReadAll(clone.Body)
		if err == nil {
			fmt.Println("[Consumer] Body", string(body))
		} else {
			fmt.Println("[Consumer] No body", err)
		}
	}
	fmt.Println("[Consumer] Cookies", req.Cookies())
	// if curl, err := s.Curl(*req); err == nil {
	// 	fmt.Println("[Consumer] Curl", curl)
	// }
}

// Curl sends an HTTP request and returns the equivalent cURL command.
// It takes an *http.Request as input and returns the cURL command as a string.
// If there is an error during the process, it returns an error.
func (s *HttpConsumerService) Curl(req http.Request) (string, error) {
	curl := []string{
		"-X", req.Method,
	}
	for key, value := range req.Header {
		curl = append(curl, "-H", fmt.Sprintf("\"%s: %s\"", key, value[0]))
	}
	if req.Body != nil {
		fmt.Println("[Consumer] Body exists")
		bodyLen := req.ContentLength
		if bodyLen > 0 {
			fmt.Println("[Consumer] Body length", bodyLen)
		} else {
			fmt.Println("[Consumer] Body length unknown")
		}
		clone := req.Clone(req.Context())
		body, err := io.ReadAll(clone.Body)
		if err == nil {
			bodyLen := len(body)
			fmt.Println("[Consumer] Body", string(body), bodyLen)
		} else {
			fmt.Println("[Consumer] No body", err)
		}
		curl = append(curl, "-d", "\""+string(body)+"\"")
	}
	curl = append(curl, req.URL.String())
	return strings.Join(curl, " "), nil
}

// Fetch sends an HTTP request and returns the response.
// It takes a pointer to an HttpRequest as input and returns a pointer to an HttResponse and an error.
// The HttpRequest contains information about the request, such as the method, URL, headers, parameters, and form data.
// The function validates the request, sets the headers, parameters, and form data, creates an HTTP request object,
// sends the request, handles authentication, handles the response, and returns the response along with any error that occurred.
func (s *HttpConsumerService) Fetch(req *HttpRequest) (*HttResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	for key, value := range s.Options.Headers {
		req.Headers = append(req.Headers, []string{key, value})
	}
	for key, value := range s.Options.Params {
		req.Params[key] = value
	}
	req.Cookies = append(req.Cookies, s.Options.Cookies...)
	targetUrl, _ := s.Options.Path(req.Path)
	var httpReq *http.Request
	var err error
	if req.Form.RawJson != "" {
		httpReq, err = http.NewRequest(req.Method, targetUrl.String(), strings.NewReader(req.Form.RawJson))
		if err != nil {
			return nil, err
		}
		req.SetContentType("application/json")
	} else if req.Form.RawText != "" {
		httpReq, err = http.NewRequest(req.Method, targetUrl.String(), strings.NewReader(req.Form.RawText))
		if err != nil {
			return nil, err
		}
		req.SetContentType("text/plain")
	} else if req.Form.RawXml != "" {
		httpReq, err = http.NewRequest(req.Method, targetUrl.String(), strings.NewReader(req.Form.RawXml))
		if err != nil {
			return nil, err
		}
		req.SetContentType("application/xml")
	} else if req.Form.RawHtml != "" {
		httpReq, err = http.NewRequest(req.Method, targetUrl.String(), strings.NewReader(req.Form.RawHtml))
		if err != nil {
			return nil, err
		}
		req.SetContentType("text/html")
	} else if req.Form.RawBinary != nil {
		httpReq, err = http.NewRequest(req.Method, targetUrl.String(), bytes.NewReader(req.Form.RawBinary))
		if err != nil {
			return nil, err
		}
		req.SetContentType("application/octet-stream")
	} else if data, ok := req.Form.Data.(map[string]string); ok {
		form := url.Values{}
		for key, value := range data {
			form.Add(key, value)
		}
		httpReq, err = http.NewRequest(req.Method, targetUrl.String(), strings.NewReader(form.Encode()))
		if err != nil {
			return nil, err
		}
		req.SetContentType("application/x-www-form-urlencoded")
	} else if req.Form.UrlEncoded != nil {
		var urlencodedForm bytes.Buffer
		multiPartWriter := multipart.NewWriter(&urlencodedForm)
		for key, value := range req.Form.UrlEncoded.(map[string]string) {
			// if value is a file
			if fileInfo, err := os.Stat(value); err == nil && !fileInfo.IsDir() {
				file, err := os.Open(value)
				if err != nil {
					continue
					// return nil, err
				}
				defer file.Close()
				fileWriter, err := multiPartWriter.CreateFormFile(key, file.Name())
				if err != nil {
					continue
					// return nil, err
				}
				if _, err := io.Copy(fileWriter, file); err != nil {
					continue
					// return nil, err
				}
			} else if strings.HasPrefix(value, "@") {
				file, err := os.Open(value[1:])
				if err != nil {
					_ = multiPartWriter.WriteField(key, value)
					continue
					// return nil, err
				}
				defer file.Close()
				fileWriter, err := multiPartWriter.CreateFormFile(key, file.Name())
				if err != nil {
					continue
					// return nil, err
				}
				if _, err := io.Copy(fileWriter, file); err != nil {
					continue
					// return nil, err
				}
			} else {
				_ = multiPartWriter.WriteField(key, value)
			}
			
		}
		multiPartWriter.Close()
		httpReq, err = http.NewRequest(req.Method, targetUrl.String(), &urlencodedForm)
		if err != nil {
			return nil, err
		}
		// req.SetContentType(multiPartWriter.FormDataContentType())
		req.SetContentType(multiPartWriter.FormDataContentType())
	} else {
		httpReq, err = http.NewRequest(req.Method, targetUrl.String(), nil)
		if err != nil {
			return nil, err
		}
	}

	for _, value := range req.Headers {
		httpReq.Header.Add(value[0], value[1])
	}
	query := httpReq.URL.Query()
	for key, value := range req.Params {
		query.Add(key, value)
	}
	httpReq.URL.RawQuery = query.Encode()
	for _, cookie := range req.Cookies {
		httpReq.AddCookie(&cookie)
	}
	httpReq = httpReq.WithContext(s.Options.Ctx)
	httpReq, err = s.HandleAuth(httpReq)
	if err != nil {
		return nil, err
	}
	// s.DebugRequest(httpReq)
	// fmt.Println("[Consumer] Request", httpReq)
	resp, err := s.Options.Http.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// fmt.Println("[Consumer] Response", resp)
	// s.DebugResponse(resp)
	resp, err = s.HandleResponse(resp, err)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	response := &HttResponse{
		Url:     resp.Request.URL.String(),
		Code:    resp.StatusCode,
		Status:  resp.Status,
		Success: resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusBadRequest,
		Headers: [][]string{},
		Cookies: func() []http.Cookie {
			cookies := []http.Cookie{}
			for _, cookie := range resp.Cookies() {
				cookies = append(cookies, *cookie)
			}
			return cookies
		}(),
	}
	for key, value := range resp.Header {
		response.Headers = append(response.Headers, []string{key, value[0]})
	}
	if _, err := response.Body.Write(body); err != nil {
		return nil, err
	}
	return response, nil
}

// FormRequest creates a new HTTP request with the specified path and form data.
// It returns the created HttpRequest object and an error, if any.
func (s *HttpConsumerService) Form(path string, form HttpRequestForm) (*HttResponse, error) {
	req := &HttpRequest{
		Method: "POST",
		Path:   path,
		Form:   form,
	}
	if resp, err := s.Fetch(req); err != nil {
		return nil, err
	} else {
		return resp, nil
	}
}

// Patch sends a PATCH request to the specified path with the provided body.
// It returns the HTTP response and an error if any.
func (s *HttpConsumerService) Patch(path string, body interface{}) (*HttResponse, error) {
	req := &HttpRequest{
		Method: "PATCH",
		Path:   path,
		Form: HttpRequestForm{
			Data: body,
		},
	}
	if resp, err := s.Fetch(req); err != nil {
		return nil, err
	} else {
		return resp, nil
	}
}

// Put sends a PUT request to the specified path with the provided body.
// It returns the HTTP response and an error if any.
func (s *HttpConsumerService) Put(path string, body interface{}) (*HttResponse, error) {
	req := &HttpRequest{
		Method: "PUT",
		Path:   path,
		Form: HttpRequestForm{
			Data: body,
		},
	}
	if resp, err := s.Fetch(req); err != nil {
		return nil, err
	} else {
		return resp, nil
	}
}

// Delete sends a DELETE request to the specified path and returns the HTTP response.
func (s *HttpConsumerService) Delete(path string, queryParams map[string]string) (*HttResponse, error) {
	req := &HttpRequest{
		Method: "DELETE",
		Path:   path,
		Params: queryParams,
	}
	if resp, err := s.Fetch(req); err != nil {
		return nil, err
	} else {
		return resp, nil
	}
}

// Post sends a POST request to the specified path with the given body.
// It returns the HTTP response and an error if any.
func (s *HttpConsumerService) Post(path string, body interface{}) (*HttResponse, error) {
	req := &HttpRequest{
		Method: "POST",
		Path:   path,
		Form: HttpRequestForm{
			Data: body,
		},
	}
	if resp, err := s.Fetch(req); err != nil {
		return nil, err
	} else {
		return resp, nil
	}
}

// Get sends a GET request to the specified path with the given query parameters.
// It returns the HTTP response and an error if any.
func (s *HttpConsumerService) Get(path string, queryParams map[string]string) (*HttResponse, error) {
	req := &HttpRequest{
		Method: "GET",
		Path:   path,
		Params: queryParams,
	}
	if resp, err := s.Fetch(req); err != nil {
		return nil, err
	} else {
		return resp, nil
	}
}

func NewHttpConsumerService(options *HttpConsumerOptions) *HttpConsumerService {
	if err := options.Validate(); err != nil {
		panic(err)
	}	
	for key, value := range HTTP_CONSUMER_DEFAULTS.Headers {
		options.Headers[key] = value
	}
	return &HttpConsumerService{
		Options: options,
	}
}
