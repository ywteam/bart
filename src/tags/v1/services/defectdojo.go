package ydk

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type DefectDojoServiceOptions struct {
	Url      string
	Token    string
	User     string
	Password string
	Adapter  *http.Client
}

func (o *DefectDojoServiceOptions) HasToken() bool {
	return len(o.Token) > 0
}
func (o *DefectDojoServiceOptions) HasCredential() bool {
	return len(o.User) > 0 && len(o.Password) > 0
}

type DefectDojoService struct {
	Options  DefectDojoServiceOptions
	Consumer *HttpServiceConsumer
}

func (s *DefectDojoService) HandleAuth(req *http.Request) (*http.Request, error) {
	fmt.Println("DefectDojo HandleAuth")
	if s.Options.HasToken() {
		req.Header.Set("Authorization", fmt.Sprintf("Token %s", s.Options.Token))
		return req, nil
	}
	if s.Options.HasCredential() {
		// req.SetBasicAuth(s.Options.User, s.Options.Password)
		// token, err := s.GetToken(s.Options.User, s.Options.Password)
		// if err != nil {
		// 	return req, err
		// }
		// req.Header.Set("Authorization", fmt.Sprintf("Token %s", token))
		return req, nil
	}
	return req, fmt.Errorf("missing token or credential")
}
func (s *DefectDojoService) HandleError(res *http.Response, err error) error {
	fmt.Println("DefectDojo HandleError")
	// body := new(DefectDojoErrorResponse)
	// body.Code = res.StatusCode

	body := &DefectDojoErrorResponse{
		Code: res.StatusCode,
	}
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("DefectDojoService ReadError: %d %s\n", res.StatusCode, err)
		return err
	}
	fmt.Printf("Response raw: %s\n", string(bodyBytes))
	if err := json.Unmarshal(bodyBytes, &body); err != nil {
		fmt.Printf("DefectDojoService ParseError: %d %s\n", res.StatusCode, err)
		return err
	}
	fmt.Printf("DefectDojoService HandleError: %d %s\n", res.StatusCode, body.Message)

	// if err := s.Consumer.BodyAs(res, body); err != nil {
	// 	fmt.Printf("DefectDojoService HandleError: %d %s\n", res.StatusCode, err)
	// 	// return err
	// }
	// fmt.Printf("DefectDojoService fetch Error: %d %s\n", body.Code, body.Message)
	// return fmt.Errorf("failed to fetch data: %s", body.Message)
	return fmt.Errorf("failed to fetch data: %s", body.Message)
}
func (s *DefectDojoService) GetToken(user string, password string) (string, error) {
	path := fmt.Sprintf("%s/api-token-auth/", s.Consumer.BaseUrl.String())
	req, err := s.Consumer.NewFormRequest(path, map[string]string{
		"username": user,
		"password": password,
	})
	// req, err := s.Consumer.NewFormRequest(path, struct {
	// 	Username string `json:"username,omitempty"`
	// 	Password string `json:"password,omitempty"`
	// }{
	// 	Username: user,
	// 	Password: password,
	// })
	if err != nil {
		return "", err
	}
	// req.SetBasicAuth(user, password)
	res, err := s.Consumer.Fetch(req)
	if err != nil {
		return "", err
	}
	body := struct {
		Token string `json:"token,omitempty"`
	}{}
	if err := s.Consumer.BodyAs(res, &body); err != nil {
		return "", err
	}
	s.Options.Token = body.Token
	return body.Token, nil
}

// func (s *DefectDojoService) Fetch(req *http.Request) (*http.Response, error) {
// 	return s.HttpClient.Fetch(req)
// }

type DefectDojoErrorResponse struct {
	Code        int      `json:"code,omitempty"`
	Detail      string   `json:"detail,omitempty"`
	Description []string `json:"description,omitempty"`
	Message     string   `json:"message,omitempty"`
}
type DefectDojoResponse struct {
	Id   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

func NewDefectDojoService(options DefectDojoServiceOptions) *DefectDojoService {
	baseUrl, err := url.Parse(options.Url + "/api/v2")
	if err != nil {
		panic(err)
	}
	service := &DefectDojoService{
		Options: options,
		Consumer: &HttpServiceConsumer{
			BaseUrl: baseUrl,
			Headers: map[string]string{},
			Ctx:     nil,
			Http: func() *http.Client {
				if options.Adapter != nil {
					return options.Adapter
				}
				return &http.Client{}
			}(),
		},
	}
	service.Consumer.AuthHandler = service.HandleAuth
	service.Consumer.ErrorHandler = service.HandleError
	return service
}
