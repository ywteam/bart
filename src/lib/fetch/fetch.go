package fetch

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"

	"yellowteam/lib/logger"
)

const (
	HTTP_CONSUMER_DEFAULT_API_KEY = "LSj!U0z^ps&Z!78Tak@f4F*#pVU8$8&0$RH!f4d#x9MnF63T5Zpz5PB6yM4M"
)

// #region Options
type Options struct {
	Ctx     context.Context // The context for the HTTP consumer.
	BaseUrl url.URL         // The base URL for the HTTP consumer.
	// Url     url.URL           // The URL for the HTTP consumer.
	Http        *http.Client      // The HTTP client for the HTTP consumer.
	Headers     map[string]string // The headers for the HTTP consumer.
	Cookies     []http.Cookie     // The cookies for the HTTP consumer.
	Params      map[string]string // The parameters for the HTTP consumer.
	Timeout     time.Duration     // The timeout for the HTTP consumer.
	Credentials ICredentailStore
}

func (o *Options) WithContext(ctx context.Context) *Options {
	o.Ctx = ctx
	return o
}
func (o *Options) WithBaseUrl(baseUrl url.URL) *Options {
	o.BaseUrl = baseUrl
	return o
}
func (o *Options) WithHttp(http *http.Client) *Options {
	o.Http = http
	return o
}
func (o *Options) WithHeaders(headers map[string]string) *Options {
	o.Headers = headers
	return o
}
func (o *Options) SetHeader(key string, value string) *Options {
	o.Headers[key] = value
	return o
}
func (o *Options) WithCookies(cookies []http.Cookie) *Options {
	o.Cookies = cookies
	return o
}
func (o *Options) WithCookie(cookie http.Cookie) *Options {
	o.Cookies = append(o.Cookies, cookie)
	return o
}
func (o *Options) WithParams(params map[string]string) *Options {
	o.Params = params
	return o
}
func (o *Options) WithParam(key string, value string) *Options {
	o.Params[key] = value
	return o
}
func (o *Options) WithTimeout(timeout time.Duration) *Options {
	o.Timeout = timeout
	return o
}
func (o *Options) WithCredentials(credentials ICredentailStore) *Options {
	o.Credentials = credentials
	return o
}
func NewClientOptions(baseUrl string, credentials ICredentailStore) *Options {
	if baseURL, err := url.Parse(baseUrl); err != nil {
		panic(err)
	} else {
		return &Options{
			Ctx:     context.Background(),
			BaseUrl: *baseURL,
			Http: &http.Client{
				Timeout: 30 * time.Second,
			},
			Headers:     map[string]string{},
			Cookies:     []http.Cookie{},
			Params:      map[string]string{},
			Timeout:     30 * time.Second,
			Credentials: credentials,
		}
	}
}

// #endregion

type FetchBody struct {
	RawText        []byte
	RawJson        []byte
	RawXml         []byte
	Form           url.Values
	FormUrlEncoded map[string]string
	FormData       map[string]string
	File           []string
}

// #region Request
type FetchRequest struct {
	Method  string
	Path    string
	Headers [][]string
	Cookies []*http.Cookie
	Params  map[string]string
	Body    FetchBody
}

func (r *FetchRequest) WithMethod(method string) *FetchRequest {
	r.Method = method
	return r
}
func (r *FetchRequest) WithPath(path string) *FetchRequest {
	r.Path = path
	return r
}
func (r *FetchRequest) WithHeaders(headers [][]string) *FetchRequest {
	r.Headers = headers
	return r
}
func (r *FetchRequest) WithCookies(cookies []*http.Cookie) *FetchRequest {
	r.Cookies = cookies
	return r
}
func (r *FetchRequest) WithParams(params map[string]string) *FetchRequest {
	r.Params = params
	return r
}
func (r *FetchRequest) WithBody(body FetchBody) *FetchRequest {
	r.Body = body
	return r
}
func (r *FetchRequest) SetHeader(key string, value string) *FetchRequest {
	r.Headers = append(r.Headers, []string{key, value})
	return r
}
func (r *FetchRequest) SetCookie(cookie *http.Cookie) *FetchRequest {
	r.Cookies = append(r.Cookies, cookie)
	return r
}
func (r *FetchRequest) SetParam(key string, value string) *FetchRequest {
	r.Params[key] = value
	return r
}
func (r *FetchRequest) SetContentType(contentType string) {
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
func (r *FetchRequest) GetContentType() string {
	for _, header := range r.Headers {
		if header[0] == "Content-Type" {
			return header[1]
		}
	}
	if len(r.Body.RawText) > 0 {
		return "text/plain"
	}
	if len(r.Body.RawJson) > 0 {
		return "application/json"
	}
	if len(r.Body.RawXml) > 0 {
		return "application/xml"
	}
	if len(r.Body.FormUrlEncoded) > 0 {
		return "application/x-www-form-urlencoded"
	}
	if len(r.Body.FormData) > 0 {
		return "multipart/form-data"
	}
	if len(r.Body.File) > 0 {
		return "application/octet-stream"
	}
	return "application/json"
}
func (r *FetchRequest) GetHeader(key string) (string, bool) {
	for _, header := range r.Headers {
		if header[0] == key {
			return header[1], true
		}
	}
	return "", false
}
func (req *FetchRequest) Make(_url *url.URL) (*http.Request, error) {
	var _request *http.Request
	var _error error
	if len(req.Body.RawText) > 0 {
		_request, _error = http.NewRequest(req.Method, _url.String(), bytes.NewBuffer(req.Body.RawText))
	} else if len(req.Body.RawJson) > 0 {
		_request, _error = http.NewRequest(req.Method, _url.String(), bytes.NewBuffer(req.Body.RawJson))
	} else if len(req.Body.RawXml) > 0 {
		_request, _error = http.NewRequest(req.Method, _url.String(), bytes.NewBuffer(req.Body.RawXml))
	} else if len(req.Body.FormUrlEncoded) > 0 {
		if req.Body.Form == nil {
			req.Body.Form = url.Values{}
		}
		for key, value := range req.Body.FormUrlEncoded {
			req.Body.Form.Add(key, value)
		}
		_request, _error = http.NewRequest(req.Method, _url.String(), bytes.NewBufferString(req.Body.Form.Encode()))
	} else if len(req.Body.FormData) > 0 {
		if req.Body.Form == nil {
			req.Body.Form = url.Values{}
		}
		for key, value := range req.Body.FormData {
			req.Body.Form.Add(key, value)
		}
		_request, _error = http.NewRequest(req.Method, _url.String(), bytes.NewBufferString(req.Body.Form.Encode()))
	} else if len(req.Body.File) > 0 {
		_request, _error = http.NewRequest(req.Method, _url.String(), bytes.NewBufferString(req.Body.File[0]))
	} else {
		_request, _error = http.NewRequest(req.Method, _url.String(), nil)
	}
	return _request, _error
}
func NewFetchRequest(method string, path string) *FetchRequest {
	return &FetchRequest{
		Method:  method,
		Path:    path,
		Headers: [][]string{},
		Cookies: []*http.Cookie{},
		Params:  map[string]string{},
		Body:    FetchBody{},
	}
}

// #endregion

// #region Response
type FetchResponse struct {
	Url     string        // The URL of the HTTP request.
	Code    int           // The HTTP status code of the response.
	Status  string        // The HTTP status message of the response.
	Success bool          // Indicates whether the request was successful or not.
	Headers [][]string    // The headers of the HTTP response.
	Cookies []http.Cookie // The cookies received in the response.
	Body    bytes.Buffer  // The body of the HTTP response.
}

func (r *FetchResponse) ExpectStatus(status ...int) bool {
	for _, s := range status {
		if r.Code == s {
			return true
		}
	}
	return false
}
func (r *FetchResponse) ExpectHeader(key string, value string) bool {
	for _, header := range r.Headers {
		if header[0] == key && header[1] == value {
			return true
		}
	}
	return false
}
func (r *FetchResponse) ExpectHeaders(headers ...string) bool {
	for _, header := range headers {
		found := false
		for _, h := range r.Headers {
			if h[0] == header {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
func (r *FetchResponse) ExpectCookie(key string, value string) bool {
	for _, cookie := range r.Cookies {
		if cookie.Name == key && cookie.Value == value {
			return true
		}
	}
	return false
}
func (r *FetchResponse) Json(v interface{}) error {
	decoder := json.NewDecoder(&r.Body)
	return decoder.Decode(&v)
}
func (r *FetchResponse) Text() string {
	return r.Body.String()
}
func (r *FetchResponse) Header(key string) (string, bool) {
	for _, header := range r.Headers {
		if header[0] == key {
			return header[1], true
		}
	}
	return "", false
}
func (r *FetchResponse) Cookie(key string) (string, bool) {
	for _, cookie := range r.Cookies {
		if cookie.Name == key {
			return cookie.Value, true
		}
	}
	return "", false
}
func NewFetchResponse(response *http.Response) *FetchResponse {
	defer response.Body.Close()
	headers := [][]string{}
	for key, values := range response.Header {
		for _, value := range values {
			headers = append(headers, []string{key, value})
		}
	}
	resp := &FetchResponse{
		Url:     response.Request.URL.String(),
		Code:    response.StatusCode,
		Status:  response.Status,
		Headers: headers,
		Cookies: func() []http.Cookie {
			cookies := []http.Cookie{}
			for _, cookie := range response.Cookies() {
				cookies = append(cookies, *cookie)
			}
			return cookies
		}(),
	}
	if response.Body != nil {
		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			logger.Error("Error reading response body %s", err.Error())
		} else {
			resp.Body = *bytes.NewBuffer(bodyBytes)
		}
	}
	return resp
}

// #endregion
type FetchClient struct {
	Options        *Options                                            // Options contains the configuration options for the HTTP consumer service.
	AuthHandler    func(*http.Request) (*http.Request, error)          // AuthHandler is a function that handles authentication for incoming requests.
	ResponseHander func(*http.Response, error) (*http.Response, error) // ResponseHandler is a function that handles the response from the HTTP request.
}

func (c *FetchClient) WithAuthHandler(handler func(*http.Request) (*http.Request, error)) *FetchClient {
	c.AuthHandler = handler
	return c
}
func (c *FetchClient) WithResponseHandler(handler func(*http.Response, error) (*http.Response, error)) *FetchClient {
	c.ResponseHander = handler
	return c
}
func (c *FetchClient) HandleAuth(req *http.Request) (*http.Request, error) {
	if c.AuthHandler != nil {
		if c.Options.Credentials != nil && !c.Options.Credentials.IsEmpty() {
			credential := c.Options.Credentials.Next()
			header := credential.Header()
			req.Header.Set(header[0], header[1])
		}
		return c.AuthHandler(req)
	}
	return req, nil
}
func (c *FetchClient) HandleResponse(res *http.Response, err error) (*http.Response, error) {
	if c.ResponseHander != nil {
		return c.ResponseHander(res, err)
	}
	return res, err
}
func (c *FetchClient) MakeRequest(req *FetchRequest) (*http.Request, error) {
	for key, value := range c.Options.Params {
		req.Path = req.Path + "?" + key + "=" + value
	}
	for key, value := range c.Options.Headers {
		req.SetHeader(key, value)
	}
	for _, cookie := range c.Options.Cookies {
		req.SetCookie(&cookie)
	}
	for key, value := range req.Params {
		req.Path = req.Path + "?" + key + "=" + value
	}
	if contentType, exists := req.GetHeader("Content-Type"); !exists || contentType == "" {
		req.SetHeader("Content-Type", req.GetContentType())
	}
	for _, value := range req.Headers {
		req.SetHeader(value[0], value[1])
	}
	for _, cookie := range req.Cookies {
		req.SetCookie(cookie)
	}
	var _error error
	_url, _error := url.Parse(c.Options.BaseUrl.String() + req.Path)
	CheckError(_error, "Error parsing URL")
	_request, _error := req.Make(_url)
	CheckError(_error, "Error creating request")
	_request = _request.WithContext(c.Options.Ctx)
	for _, header := range req.Headers {
		_request.Header.Set(header[0], header[1])
	}
	for _, cookie := range req.Cookies {
		_request.AddCookie(cookie)
	}
	return _request, nil
}
func (c *FetchClient) Fetch(req *FetchRequest) (*FetchResponse, error) {
	_request, err := c.MakeRequest(req)
	CheckError(err, "Error making request")
	_request, err = c.HandleAuth(_request)
	CheckError(err, "Error handling authentication")
	_response, err := c.Options.Http.Do(_request)
	CheckError(err, "Error making request")
	_response, err = c.HandleResponse(_response, err)
	CheckError(err, "Error handling response")
	return NewFetchResponse(_response), nil
}
func (c *FetchClient) Get(path string) (*FetchResponse, error) {
	return c.Fetch(NewFetchRequest("GET", path))
}
func (c *FetchClient) Post(path string, body interface{}) (*FetchResponse, error) {
	return c.Fetch(NewFetchRequest("POST", path).WithBody(FetchBody{
		RawJson: func() []byte {
			if b, err := json.Marshal(body); err != nil {
				logger.Error("Error marshalling body %s", err.Error())
				return []byte{}
			} else {
				return b
			}
		}(),
	}))
}
func (c *FetchClient) Put(path string) (*FetchResponse, error) {
	return c.Fetch(NewFetchRequest("PUT", path))
}
func (c *FetchClient) Patch(path string) (*FetchResponse, error) {
	return c.Fetch(NewFetchRequest("PATCH", path))
}
func (c *FetchClient) Delete(path string) (*FetchResponse, error) {
	return c.Fetch(NewFetchRequest("DELETE", path))
}
func (c *FetchClient) Head(path string) (*FetchResponse, error) {
	return c.Fetch(NewFetchRequest("HEAD", path))
}

func CheckError(err error, message ...string) {
	if err != nil {
		logger.Error(err.Error(), message)
	}
}
func NewClient(options *Options) *FetchClient {
	options.WithContext(context.Background())
	return &FetchClient{
		Options: options,
	}
}
