package defectdojo

import (
	// "fmt"
	// "net/url"
	"encoding/json"
	"fmt"
	fetch "github.com/ywteam/ydk-go/lib/fetch"
	"net/http"
	"strings"
)

type DefectDojoServiceBasicAuth struct {
	Username string
	Password string
}
type DefectDojoServiceTokenAuth struct {
	Token string
}

type DefectDojoServiceOptions struct {
	Url  string
	Auth interface{}
}

type DefectDojoService struct {
	Options DefectDojoServiceOptions
	Client  *fetch.FetchClient
}

func (s *DefectDojoService) Auth(credential interface{}) (string, error) {
	switch v := credential.(type) {
	case DefectDojoServiceBasicAuth:
		if token, err := s.GetToken(v.Username, v.Password); err != nil {
			return "", err
		} else {
			// s.Client.Options.Credentials.AddBearerToken(token)
			return token, nil
		}
	case DefectDojoServiceTokenAuth:
		// fmt.Println("Token Auth", v.Token)
		// s.Client.Options.Credentials.AddBearerToken(v.Token)
		return v.Token, nil
	default:
		fmt.Println("Unknown Auth")
	}
	return "", nil
}
func (s *DefectDojoService) GetToken(username, password string) (string, error) {
	res, err := s.Client.Post("/api-token-auth/", map[string]string{
		"username": username,
		"password": password,
	})
	if err != nil {
		return "", err
	}
	if res.Code != 200 {
		return "", fmt.Errorf("invalid credentials")
	}
	var body struct {
		Token string `json:"token"`
	}
	if err := res.Json(&body); err != nil {
		return "", err
	}
	return body.Token, nil
}
func (s *DefectDojoService) HandleAuth(req *http.Request) (*http.Request, error) {
	if strings.Contains(req.URL.Path, "/api-token-auth/") {
		return req, nil
	}
	if s.Client.Options.Credentials.IsEmpty() {
		if token, err := s.Auth(s.Options.Auth); err != nil {
			return nil, err
		} else {
			s.Client.Options.Credentials.AddBearerToken(token)
		}
	}
	return req, nil
}
func (s *DefectDojoService) HandleResponse(res *http.Response, err error) (*http.Response, error) {
	if err != nil {
		if res.StatusCode >= 300 {
			body := DefectDojoErrorResponse{
				Code: res.StatusCode,
			}
			if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
				return res, err
			}

			if res.StatusCode == 401 {
				if token, err := s.Auth(s.Options.Auth); err == nil {
					s.Client.Options.Credentials.Rotate(token)
					request := res.Request.Clone(res.Request.Context())
					request.Header.Del("Authorization")
					return s.Client.Options.Http.Do(request)
					// return res, nil
				}
				return res, err
			}
			return res, fmt.Errorf("DefectDojo Error: %s", body.Message)
		}
		return res, err
	}

	return res, err
}
// func (s *DefectDojoService) Next(res *DefectDojoResponse) (*DefectDojoResponse, error) {
// 	if res.Next != "" {
// 		url = res.
// 		resp, err := s.Client.Get(res.Next)
// 		if err != nil {
// 			return nil, err
// 		}
// 		var response DefectDojoResponse
// 		if err := resp.Json(&response); err != nil {
// 			return nil, err
// 		}
// 		return &response, nil
// 	}
// 	return nil, nil
// }
func NewDefectDojoService(options DefectDojoServiceOptions) (*DefectDojoService, error) {
	client := fetch.NewClient(fetch.NewClientOptions(options.Url+"/api/v2", fetch.NewCredentialStore()))
	return &DefectDojoService{
		Options: options,
		Client:  client,
	}, nil
}

// func NewDefectDojoServiceV1(url string, args ...interface{}) (*DefectDojoService, error) {
// 	options := DefectDojoServiceOptions{
// 		Url: url,
// 	}
// 	for _, arg := range args {
// 		switch v := arg.(type) {
// 		case DefectDojoServiceOptions:
// 			options = v
// 		case DefectDojoServiceBasicAuth:
// 			options.Auth = v
// 		case DefectDojoServiceTokenAuth:
// 			options.Auth = v
// 		default:
// 			return nil, fmt.Errorf("unsupported argument type: %T", v)
// 		}
// 	}
// 	client := fetch.NewClient(fetch.NewClientOptions(options.Url+"/api/v2", fetch.NewCredentialStore()))
// 	return &DefectDojoService{
// 		Options: options,
// 		Client:  client,
// 	}, nil
// }

// func NewDefectDojoService(options DefectDojoServiceOptions) (*DefectDojoService, error) {
// 	if err := options.Validate(); err != nil {
// 		return nil, err
// 	}
// 	clientOptions := fetch.NewClientOptions(options.Url + "/api/v2")
// 	clientOptions.WithHeaders(map[string]string{
// 		"Content-Type": "application/json",
// 	})
// 	client := fetch.NewClient(clientOptions)
// 	return &DefectDojoService{
// 		Options: options,
// 		Client:  client,
// 	}, nil
// }

func init() {
	store := fetch.NewCredentialStore() //.AddBasicAuth("admin", "admin").AddBearerToken("token")
	store.AddBasicAuth("admin", "admin")
	store.AddBearerToken("token")
	store.AddCustomHeader("X-Api-Key", "key")
	options := DefectDojoServiceOptions{}
	client := fetch.NewClient(fetch.NewClientOptions("http://localhost:8000/api/v2", store))
	service := DefectDojoService{
		Options: options,
		Client:  client,
	}
	fmt.Println(service)
}
