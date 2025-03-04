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

type HttpConsumer struct {
	Ctx          context.Context
	BaseUrl      string
	Http         *http.Client
	Headers      map[string]string
	Cookies      []*http.Cookie
	AuthHandler  func(*http.Request) (*http.Request, error)
	ErrorHandler func(*http.Response, error) error
}
func (c *HttpConsumer) Curl(req *http.Request) (string, error) {
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
func (c *HttpConsumer) Fetch(req *http.Request) (*http.Response, error) {
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
	url, err := url.Parse(c.BaseUrl)
	if err != nil {
		return nil, err
	}
	req.URL = url.ResolveReference(req.URL)		
	fmt.Println("Http Consumer Fetch", req.URL)
	bodyBytes, _ := io.ReadAll(req.Body)
	fmt.Printf("Http Consumer FetchBody: %s\n", bodyBytes)
	res, err := c.Http.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	bodyBytes, err = io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}	
	fmt.Printf("Http Consumer FetchResponse: %s, %s\n", res.Status, bodyBytes)
	curl, _ := c.Curl(req)
	fmt.Println("Curl: " + curl)
	return res, nil
	// return c.Http.Do(req)
}
func (c *HttpConsumer) NewRequest(method string, path string) (*http.Request, error) {
	fmt.Println("HttpConsumer NewRequest")
	if c.Ctx == nil {
		c.Ctx = context.Background()
	}
	if c.Http == nil {
		c.Http = &http.Client{}
	}
	if c.Headers == nil {
		c.Headers = make(map[string]string)
	}
	if c.Cookies == nil {
		c.Cookies = make([]*http.Cookie, 0)
	}
	return http.NewRequestWithContext(c.Ctx, method, c.BaseUrl+path, nil)
}
func (c *HttpConsumer) NewRequestWithBody(method string, path string, body interface{}) (*http.Request, error) {
	fmt.Println("HttpConsumer NewRequestWithBody")
	request, err := c.NewRequest(method, path)
	if err != nil {
		return nil, err
	}
	if body != nil {
		request.Header.Add("Content-Type", "application/json")
		json, err := json.Marshal(body)
		if err == nil {
			request.Body = io.NopCloser(bytes.NewReader(json))
		} else {
			switch v := body.(type) {
			case io.Reader:
				request.Body = io.NopCloser(v)
			case []byte:
				request.Body = io.NopCloser(bytes.NewReader(v))
			case string:
				request.Body = io.NopCloser(strings.NewReader(v))
			default:
				return nil, fmt.Errorf("unsupported body type: %T", body)
			}
		}
	}
	return request, nil
}
func (c *HttpConsumer) NewRequestWithForm(method string, path string, form interface{}) (*http.Request, error) {
	fmt.Println("HttpConsumer NewRequestWithForm")
	request, err := c.NewRequest(method, path)
	if err != nil {
		return nil, err
	}
	if form != nil {
		request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		values := url.Values{} //request.URL.Query()
		for key, value := range form.(map[string]string) {
			values.Add(key, value)
		}
		request.Body = io.NopCloser(strings.NewReader(values.Encode()))		
	}
	rawBody, _ := io.ReadAll(request.Body)
	fmt.Printf("NewRequestWithForm: HttpConsumer NewRequestWithForm body: %s\n", rawBody)
	return request, nil

}

type HttpExampleService struct {
	Consumer *HttpConsumer
}

func (s *HttpExampleService) GetExample() (*http.Response, error) {
	req, err := s.Consumer.NewRequestWithForm(http.MethodPost, "/api-token-auth", map[string]string{
		"username": "admin",
		"password": "1Defectdojo@demo#appsec",
	})
	if err != nil {
		return nil, err
	}
	bodyBytes, _ := io.ReadAll(req.Body)
	fmt.Printf("GetExample: Http Consumer FetchBody: %s\n", bodyBytes)
	fmt.Println(req.Method + ": HttpExampleService GetExample: " + req.URL.String())	
	res, err := s.Consumer.Fetch(req)
	if err != nil {
		return nil, err
	}
	return res, nil
	// payload := strings.NewReader("username=admin&password=1Defectdojo%40demo%23appsec")
	// req, err := http.NewRequest(http.MethodPost, "/api-token-auth", payload)
	// if err != nil {
	// 	return nil, err
	// }
	// fmt.Println("HttpExampleService GetExample: " + req.URL.String())
	// req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	// req.Header.Add("Authorization", "Basic YWRtaW46MURlZmVjdGRvam9AZGVtbyNhcHBzZWM=")
	// res, err := s.Consumer.Fetch(req)
	// if err != nil {
	// 	return nil, err
	// }
	// return res, nil
}
func NewHttpExampleService() *HttpExampleService {
	return &HttpExampleService{
		Consumer: &HttpConsumer{
			BaseUrl: "https://demo.defectdojo.org/api/v2",
			Http:    &http.Client{},
		},
	}
}

// func HttpConsumerAuth() {

// 	url := "https://demo.defectdojo.org/api/v2/api-token-auth/"
// 	method := "POST"

// 	payload := strings.NewReader("username=admin&password=1Defectdojo%40demo%23appsec")

// 	client := &http.Client{}
// 	req, err := http.NewRequest(method, url, payload)

// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
// 	req.Header.Add("Authorization", "Basic YWRtaW46MURlZmVjdGRvam9AZGVtbyNhcHBzZWM=")

// 	res, err := client.Do(req)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	defer res.Body.Close()

// 	body, err := io.ReadAll(res.Body)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	fmt.Println(string(body))
// }
