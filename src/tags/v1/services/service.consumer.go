package ydk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

var DEFAULT_HEADERS = map[string]string{
	"User-Agent": "YDK/1.0",
	// "Accept": "application/json",
	// "Content-Type": "application/json",
	"Y-Api-Key": "1234567890",
	"Y-Version": "1.0",
}

type HttpServiceConsumer struct {
	Ctx          context.Context
	BaseUrl      *url.URL
	Http         *http.Client
	Headers      map[string]string
	Cookies      []*http.Cookie
	AuthHandler  func(*http.Request) (*http.Request, error)
	ErrorHandler func(*http.Response, error) error
}

func (c *HttpServiceConsumer) HandleAuth(req *http.Request) (*http.Request, error) {
	fmt.Println("Http Consumer  HandleAuth")
	if c.AuthHandler != nil {
		return c.AuthHandler(req)
	}
	return req, nil
}
func (c *HttpServiceConsumer) HandleError(res *http.Response, err error) error {
	fmt.Println("Http Consumer HandleError")
	if c.ErrorHandler != nil {
		return c.ErrorHandler(res, err)
	}
	fmt.Println("Http Consumer HandleError done")
	return err
}
func (c *HttpServiceConsumer) BodyAs(res *http.Response, v interface{}) error {
	if err := json.NewDecoder(res.Body).Decode(v); err != nil {
		fmt.Printf("Http Consumer BodyAs: %v\n", err)
		return err
	}
	return nil
}
func (c *HttpServiceConsumer) NewRequest(method string, path string, body interface{}) (*http.Request, error) {
	var req *http.Request
	var err error
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		fmt.Printf("NewRequest body %v \n", string(b))
		req, err = http.NewRequest(method, path, bytes.NewReader(b))
		if err != nil {
			return nil, err
		}
	} else {
		req, err = http.NewRequest(method, path, nil)
		if err != nil {
			return nil, err
		}
	}

	return req, nil
}
func (c *HttpServiceConsumer) NewFormRequest(path string, body interface{}) (*http.Request, error) {
	req, err := c.NewRequest(http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}
	// check if body has username and password
	if body != nil {
		if v, ok := body.(map[string]string); ok {
			// print user name
			if user, ok := v["username"]; ok {
				fmt.Printf("NewFormRequest username %v \n", user)
			}
			if password, ok := v["password"]; ok {
				fmt.Printf("NewFormRequest password %v \n", password)
			}
			form := url.Values{}
			for key, value := range v {
				form.Add(key, value)
			}
			req.Body = io.NopCloser(strings.NewReader(form.Encode()))
		}
	}

	fmt.Printf("NewFormRequest body %v \n", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req, nil

}
func (c *HttpServiceConsumer) NewPostRequest(path string, body interface{}) (*http.Request, error) {
	return c.NewRequest(http.MethodPost, path, body)
}
func (c *HttpServiceConsumer) NewPutRequest(path string, body interface{}) (*http.Request, error) {
	return c.NewRequest(http.MethodPut, path, body)
}
func (c *HttpServiceConsumer) NewPatchRequest(path string, body interface{}) (*http.Request, error) {
	return c.NewRequest(http.MethodPatch, path, body)
}
func (c *HttpServiceConsumer) NewDeleteRequest(path string) (*http.Request, error) {
	return c.NewRequest(http.MethodDelete, path, nil)
}
func (c *HttpServiceConsumer) NewGetRequest(path string) (*http.Request, error) {
	return c.NewRequest(http.MethodGet, path, nil)
}
func (c *HttpServiceConsumer) Curl(req *http.Request) (string, error) {
	curl := []string{
		"-X", req.Method,
	}
	for key, value := range req.Header {
		curl = append(curl, "-H", fmt.Sprintf("\"%s: %s\"", key, value[0]))
	}
	if req.Body != nil {
		b, err := io.ReadAll(req.Body)
		if err != nil {
			return "", err
		}
		curl = append(curl, "-d", "\""+string(b)+"\"")
	}
	curl = append(curl, req.URL.String())
	return strings.Join(curl, " "), nil

}
func (c *HttpServiceConsumer) Debug(req *http.Request) {
	curl, err := c.Curl(req)
	if err != nil {
		fmt.Printf("Http Consumer DebugError: %s\n", err)
		return
	}
	fmt.Printf("Http Consumer Debug: %s\n", curl)
}
func (c *HttpServiceConsumer) Fetch(req *http.Request) (*http.Response, error) {
	if c.Ctx != nil {
		req = req.WithContext(c.Ctx)
	}
	for key, value := range DEFAULT_HEADERS {
		if _, ok := c.Headers[key]; !ok {
			c.Headers[key] = value
		}
	}
	for key, value := range c.Headers {
		if len(req.Header.Get(key)) == 0 {
			req.Header.Set(key, value)
		}
	}
	for _, cookie := range c.Cookies {
		req.AddCookie(cookie)
	}
	c.BaseUrl.Path = req.URL.Path
	req.URL = c.BaseUrl
	fmt.Println("Http Consumer Fetch", req.URL)
	// req, err := c.HandleAuth(req)
	// if err != nil {
	// 	return nil, err
	// }
	c.Debug(req)
	res, err := c.Http.Do(req)
	// res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Http Consumer FetchError: %s\n", err)
		return nil, err
	}
	fmt.Printf("Http Consumer FetchResponse: %s\n", res.Status)
	// defer res.Body.Close()
	defer func() {
		fmt.Println("Http Consumer Fetch defer")
		if res.Body != nil {
			res.Body.Close()
		}
	}()
	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		fmt.Printf("Http Consumer FetchError: %d %s\n", res.StatusCode, res.Status)
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Printf("Http Consumer ReadError: %d %s\n", res.StatusCode, err)
			return nil, err
		}
		if len(bodyBytes) == 0 {
			return res, fmt.Errorf("empty body %d: %s", res.StatusCode, res.Status)
		}
		fmt.Printf("(%d) Response raw: %s (%d)\n", res.StatusCode, string(bodyBytes), len(bodyBytes))
		return res, fmt.Errorf("HTTP %d: %s", res.StatusCode, res.Status)
		// return res, c.HandleError(res, fmt.Errorf("HTTP %d: %s", res.StatusCode, res.Status))
	}
	fmt.Println("Http Consumer Fetch done")
	return res, nil
}
func (c *HttpServiceConsumer) Get(path string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	return c.Fetch(req)
}
func (c *HttpServiceConsumer) Post(path string, body interface{}) (*http.Response, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, path, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	return c.Fetch(req)
}
func (c *HttpServiceConsumer) Put(path string, body interface{}) (*http.Response, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPut, path, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	return c.Fetch(req)
}
func (c *HttpServiceConsumer) Delete(path string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}
	return c.Fetch(req)
}
func (c *HttpServiceConsumer) Patch(path string, body interface{}) (*http.Response, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPatch, path, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	return c.Fetch(req)
}

func NewHttpServiceConsumer(baseUrl string, ctx context.Context) (*HttpServiceConsumer, error) {
	u, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}
	if ctx == nil {
		ctx = context.Background()
	}
	return &HttpServiceConsumer{
		BaseUrl: u,
		Http:    http.DefaultClient,
		Headers: map[string]string{},
		Cookies: []*http.Cookie{},
		Ctx:     ctx,
	}, nil
}

// curl -X POST -H 'Content-Type: application/x-www-form-urlencoded' -H 'Authorization: Basic YWRtaW46MURlZmVjdGRvam9AZGVtbyNhcHBzZWM=' -d 'username=admin' -d 'password=1Defectdojo@demo#appsec' https://demo.defectdojo.org/api/v2/api-token-auth/
