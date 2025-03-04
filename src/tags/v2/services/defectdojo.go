package ydk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type DefectDojoServiceOptions struct {
	Url      string
	Token    string
	User     string
	Password string
}

func (o *DefectDojoServiceOptions) HasToken() bool {
	return len(o.Token) > 0
}
func (o *DefectDojoServiceOptions) HasCredential() bool {
	return len(o.User) > 0 && len(o.Password) > 0
}
func (o *DefectDojoServiceOptions) Validate() error {
	if len(o.Url) == 0 {
		return fmt.Errorf("missing URL")
	}
	if !o.HasToken() && !o.HasCredential() {
		return fmt.Errorf("missing token or credential")
	}
	if _, err := url.Parse(o.Url); err != nil {
		return err
	}
	return nil
}

type DefectDojoService struct {
	Options   DefectDojoServiceOptions
	Consumer  *HttpConsumerService
	Languages *DefectDojoLanguageService
}

func (s *DefectDojoService) HandleAuth(req *http.Request) (*http.Request, error) {
	if strings.Contains(req.URL.Path, "/api-token-auth/") {
		return req, nil
	}
	fmt.Println("DefectDojo HandleAuth")
	if s.Options.HasToken() {
		// fmt.Println("Configuring token", s.Options.Token)
		req.Header.Set("Authorization", fmt.Sprintf("Token %s", s.Options.Token))
		return req, nil
	}
	if s.Options.HasCredential() {
		// fmt.Println("Configuring credential", s.Options.User, s.Options.Password)
		req.SetBasicAuth(s.Options.User, s.Options.Password)
		return req, nil
		// if token, err := s.GetToken(s.Options.User, s.Options.Password); err == nil {
		// 	req.Header.Set("Authorization", fmt.Sprintf("Token %s", token))
		// 	return req, nil
		// } else {
		// 	return req, err
		// }
	}

	return req, fmt.Errorf("missing token or credential")
}
func (s *DefectDojoService) HandleResponse(res *http.Response, err error) (*http.Response, error) {
	if err != nil {
		if res.StatusCode <= 200 || res.StatusCode >= 400 {
			body := DefectDojoErrorResponse{
				Code: res.StatusCode,
			}
			if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
				return res, err
			}
			return res, fmt.Errorf("DefectDojo Error: %s", body.Message)
		}
		return res, err
	}

	return res, err
}
func (s *DefectDojoService) GetToken(user string, password string) (string, error) {
	if s.Options.HasToken() {
		return s.Options.Token, nil
	}
	payload := map[string]string{
		"username": user,
		"password": password,
	}
	res, err := s.Consumer.Form("/api-token-auth/", HttpRequestForm{
		UrlEncoded: payload,
	})
	if err != nil {
		return "", err
	}
	fmt.Println("GetToken", res.Text())

	var body struct {
		Token string `json:"token"`
	}
	if err := res.Json(&body); err != nil {
		fmt.Println("Error", err)
		return "", err
	}
	s.Options.Token = body.Token
	return body.Token, nil
}
func NewDefectDojoService(options DefectDojoServiceOptions) *DefectDojoService {
	options.Validate()
	service := &DefectDojoService{
		Options: options,
		Consumer: NewHttpConsumerService(&HttpConsumerOptions{
			Ctx:     HTTP_CONSUMER_DEFAULTS.Ctx,
			BaseUrl: options.Url + "/api/v2",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}),
	}
	service.Consumer.AuthHandler = service.HandleAuth
	service.Consumer.ResponseHander = service.HandleResponse
	if !options.HasToken() {
		if _, err := service.GetToken(options.User, options.Password); err != nil {
			fmt.Println("Authenication Error", err)
		}
	}
	// fmt.Println("Token configured", service.Options.Token)
	service.Languages = &DefectDojoLanguageService{Api: service}
	return service
}

type DefectDojoErrorResponse struct {
	Code        int      `json:"code,omitempty"`
	Detail      string   `json:"detail,omitempty"`
	Description []string `json:"description,omitempty"`
	Message     string   `json:"message,omitempty"`
}

type DefectDojoResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []interface{}
}
